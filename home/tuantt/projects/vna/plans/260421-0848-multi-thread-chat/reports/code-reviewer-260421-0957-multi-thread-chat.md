# Code Review — Multi-Thread Chat Feature

**Date:** 2026-04-21 | **Branch:** feat/thread-message (merged to main)
**Reviewer:** code-reviewer agent

---

## Scope

- Files: 14 (SQL migration, domain, repository, service, handler, delivery, main)
- Build: `go build ./...` ✓ | Tests: 23/23 ✓
- Focus: security, correctness, error handling, YAGNI/KISS/DRY

---

## Overall Assessment

Solid implementation. Architecture is clean, layering is well-respected, IDOR is addressed. Two correctness issues need fixing (one high, one medium). No critical security vulnerabilities.

---

## Critical Issues

None.

---

## High Priority

### 1. Orphan thread on upstream failure (service/message.go:83-149)

`resolveThread` creates a new thread **before** calling OpenResponses. If the AI call fails, the thread row persists in the DB with no messages — client receives an error but the orphan thread appears in `GET /api/v1/threads`.

```go
// Current flow in SendStream / Send:
resolvedID, err := s.resolveThread(...)   // thread row inserted here
// ... AI call fails ...
return nil, err                            // orphan thread left behind
```

**Fix options (pick one):**

- **Lazy creation (preferred):** pass `threadID` through to AI call, save thread only on successful message save. Wrap thread insert + message insert in a DB transaction.
- **Immediate:** Wrap the entire `resolveThread` + AI call + `Save` in a `sql.Tx`. Rollback on any error. Only commit when message is saved.

Neither is trivial but the lazy approach respects KISS better — defer thread creation until just before `repo.Save`.

### 2. `errors.Is(sql.ErrNoRows)` silently breaks after wrapping (repository/thread.go:75, service/message.go:94)

`GetByIDAndUser` wraps errors with `fmt.Errorf("repository.thread: get: %w", err)`. The service checks:

```go
if errors.Is(err, sql.ErrNoRows) { ... }
```

`errors.Is` does unwrap `%w` chains, so this actually works in Go 1.13+. **But** the check in `service/message.go` receives the wrapped error — verify the sentinel propagates correctly end-to-end.

**Quick test to confirm:** add a unit test asserting `errors.Is(err, sql.ErrNoRows)` after `GetByIDAndUser` returns a not-found result. Current tests mock at the repo interface level and bypass this wrapping.

_Severity reduced from Critical to High because Go's `%w` + `errors.Is` does unwrap — but this is fragile and worth an explicit test._

---

## Medium Priority

### 3. Duplicate limit defaulting (service/message.go:174, handler/http/message.go:184-190)

Both handler and service apply the same default limit=20:

- Handler: defaults `limit = 20` then validates before passing to service
- Service: `if limit == 0 { limit = 20 }`

The service branch (`limit == 0`) can never be reached when called from the handler, because the handler never passes 0 (it applies the default first). The service branch only matters if `ListByThread` is called directly (e.g. tests).

**Fix:** Remove the handler's defaulting and let the service own the business rule, or document that the handler is the canonical entry point and remove service's dead branch. Either way — one place.

### 4. No input validation on `thread_id` and `cursor` UUID format

`thread_id` in query params and `cursor` in `ListByThread` are passed directly to Postgres. Invalid UUID strings (e.g. `"abc"`) cause a Postgres `invalid input syntax for type uuid` error which propagates as a 500 Internal Server Error rather than 400.

```
GET /api/v1/messages?thread_id=notauuid  → 500 (should be 400)
GET /api/v1/messages?thread_id=...&cursor=notauuid → 500 (should be 400)
```

**Fix:** Add a lightweight UUID format check before service call:

```go
// simple check without extra dep
import "regexp"
var uuidRe = regexp.MustCompile(`(?i)^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)
```

Or use `github.com/google/uuid` which is likely already transitive.

### 5. Thread name is permanently "New Chat" with no rename capability

Schema stores `name TEXT NOT NULL DEFAULT 'New Chat'` but there is no `PATCH /api/v1/threads/:id` endpoint or service method. Fine for MVP, but the field creates expectation without fulfillment.

**Fix (YAGNI-safe):** Either drop the `name` column until rename is built, or add it to the roadmap. Keeping it as dead weight violates YAGNI.

---

## Low Priority

### 6. No `updated_at` / `last_message_at` on threads

`GET /api/v1/threads` orders by `created_at DESC`. In practice, users expect threads ordered by most recent activity. This will require a schema change later.

**Suggestion:** Add `last_message_at TIMESTAMPTZ` to threads, updated via trigger or in `repo.Save`. Not urgent but cheaper now than as a migration later.

### 7. `listResponse.Data` returns `null` instead of `[]` on empty (repository/message.go:103)

When `msgs` is empty, `ListByThread` returns the nil slice. In `handler/http/message.go`:

```go
data := make([]*dto.MessageResponse, len(msgs))  // len=0 → non-nil empty slice ✓
```

Handler correctly uses `make`, so JSON output is `[]` not `null`. This is fine — noted for awareness only.

### 8. 5-minute SSE timeout is hardcoded (handler/http/message.go:123)

```go
ctx, cancel := context.WithTimeout(c.Request().Context(), 5*time.Minute)
```

Consider pulling this from config for tuning without redeployment.

---

## Edge Cases Found by Scout

1. **Orphan thread on AI failure** — covered in High #1 above
2. **Invalid UUID inputs returning 500** — covered in Medium #4 above
3. **`errors.Is` wrapping correctness** — covered in High #2 above
4. **Duplicate limit default logic** — covered in Medium #3 above
5. **No thread activity ordering** — covered in Low #6 above

---

## Positive Observations

- IDOR protection solid: `GetByIDAndUser` enforces `WHERE id=$1 AND user_id=$2` at SQL level; ownership re-validated in both `resolveThread` and `ListByThread`
- SSE metadata-first guarantee: `onMeta` called synchronously before `streamOpenResponses` — ordering is deterministic
- Keyset pagination implementation is correct and uses compound `(created_at, id)` tie-breaking
- Error sentinel design is clean; all service errors are typed and correctly mapped to HTTP codes
- 1MB scanner buffer for SSE stream is a thoughtful guard against large delta payloads
- `sql.ErrNoRows` mapped to 404 not 403 — correct for IDOR (not revealing thread existence to unauthorized users)
- Logging uses structured `slog` throughout; no sensitive data (question/answer text) logged at error level

---

## Recommended Actions

1. **[High]** Fix orphan thread: wrap `resolveThread` + AI call + `repo.Save` in a transaction, or defer thread creation until save
2. **[High]** Add an explicit integration-style test for `errors.Is(err, sql.ErrNoRows)` propagation through `GetByIDAndUser` wrapper
3. **[Medium]** Validate UUID format for `thread_id` and `cursor` query params; return 400 for malformed values
4. **[Medium]** Remove duplicate limit defaulting — single source of truth in service layer
5. **[Low]** Decide on `name` column: implement rename endpoint or remove the column (YAGNI)
6. **[Low]** Move SSE timeout to config

---

## Metrics

- Type Coverage: 100% (Go is statically typed; no `interface{}` misuse found)
- Test Coverage: handler=52.4%, service=7.1% (low service coverage — mocks bypass real path)
- Linting Issues: 0 reported by `go build`
- Security: No SQL injection risk (parameterized queries throughout), no token/secret leakage in logs

---

## Unresolved Questions

1. Should `POST /api/v1/messages` (non-streaming) also support implicit thread creation on empty `thread_id`? Currently it does (same `resolveThread`), but the design doc only mentions the stream endpoint. Confirm intentional.
2. Is there a plan for thread rename (`PATCH /api/v1/threads/:id`)? If yes, the `name` column is justified; if no timeline, consider removing it.
3. OpenResponses `user` param = `thread_id` — does this correctly isolate AI context per thread in the upstream provider, or should it be `user_id:thread_id` to avoid cross-user collisions if thread IDs are ever shared?
