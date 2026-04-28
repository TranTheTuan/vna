# Phase 06 — Tests

**Context:** `plans/reports/brainstorm-260421-0848-multi-thread-chat.md`
**Priority:** P2 | **Status:** Pending | **Effort:** 1h
**Blocked by:** Phase 05

## Overview

Check existing test files, extend/add tests for thread creation, thread-scoped message flow, and `event: metadata` SSE emission.

## Related Code Files

| Action | File |
|--------|------|
| Explore | existing `*_test.go` files |
| Modify/Create | test files for message and thread layers |

## Implementation Steps

### 1. Explore existing tests

```bash
find /home/tuantt/projects/vna -name "*_test.go" | sort
```

Understand patterns before writing new tests — follow existing style exactly.

### 2. Repository-level tests (if integration tests exist)

- `threadRepository.Create` → inserts row, returns correct fields
- `threadRepository.GetByIDAndUser` → returns `sql.ErrNoRows` for wrong user
- `messageRepository.Save` → persists `thread_id`
- `messageRepository.ListByThread` → returns only messages for that thread, correct keyset pagination

### 3. Service-level tests (unit, mock repos)

- `messageService.SendStream` with empty `threadID` → calls `threadRepo.Create`, emits `onMeta` with new ID before any `onDelta`
- `messageService.SendStream` with valid `threadID` → validates ownership, emits `onMeta` with same ID
- `messageService.SendStream` with foreign `threadID` → returns `ErrThreadNotFound`
- `messageService.ListByThread` with foreign `threadID` → returns `ErrThreadNotFound`

### 4. Handler-level tests (httptest)

- `POST /stream` with empty `thread_id` → SSE: first event is `metadata`, next are `delta`, final is `done` with `thread_id`
- `GET /threads` → 200 with list (empty or populated)
- `GET /messages?thread_id=xxx` → 400 when `thread_id` missing
- `GET /messages?thread_id=foreign` → 404

### 5. Run tests

```bash
cd /home/tuantt/projects/vna && go test ./...
```

Fix any failures before marking phase complete.

## Todo

- [x] Explore existing `*_test.go` files to understand test patterns
- [x] Add/extend repository tests for `ThreadRepository` and updated `MessageRepository`
- [x] Add service unit tests for thread resolution and `onMeta` callback ordering
- [x] Add handler tests for `event: metadata` SSE ordering and `GET /threads`
- [x] Run `go test ./...` — all pass

## Success Criteria

- `go test ./...` exits 0
- `event: metadata` is confirmed as the first SSE event in stream tests
- Thread ownership validation confirmed via test (foreign thread_id → 404/error)

## Risk Assessment

- If no integration test harness exists, focus on unit tests with mocked repos
- Do not add a test framework — use whatever is already in the codebase
