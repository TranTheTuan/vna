# Phase 02 — Wire-up & Integration

**Context:** `plans/260428-health-check-endpoints/plan.md`
**Priority:** P2 | **Status:** Complete | **Effort:** 0.2h
**Blocked by:** Phase 01

## Overview

Register health routes in `cmd/api/main.go` by passing DB pool to `RegisterHealthRoutes`.

## Requirements

- Call `http_delivery.RegisterHealthRoutes(e, pool)` after pool creation
- Use `e` (Echo instance), not `apiGroup` — health endpoints at root level
- No JWT middleware on health endpoints (already handled in phase 01)

## Related Code Files

| Action | File |
|--------|------|
| Modify | `cmd/api/main.go` |

## Implementation Steps

1. In `cmd/api/main.go`, after line 82 (after all other route registrations), add:

```go
// Health check endpoints (no JWT middleware)
http_delivery.RegisterHealthRoutes(e, pool)
```

2. Ensure `database/sql` is imported (needed for `*sql.DB` type in `RegisterHealthRoutes`):
   - Check if `github.com/TranTheTuan/vna/internal/db` already imports `database/sql`
   - If not, add `"database/sql"` to imports in `health.go` (not in `main.go`)

## Todo

- [ ] Add `http_delivery.RegisterHealthRoutes(e, pool)` call in `cmd/api/main.go`
- [ ] Verify `database/sql` import is available where needed
- [ ] Run `go build ./cmd/api/...` to verify compilation

## Success Criteria

- `go build ./cmd/api/...` succeeds
- Health routes registered at root level (not under `/api/v1`)
- `pool` (type `*sql.DB`) passed correctly to `RegisterHealthRoutes`

## Risk Assessment

- Low risk — one line addition to main.go
- Ensure import paths are correct (`http_delivery` alias matches existing code)
