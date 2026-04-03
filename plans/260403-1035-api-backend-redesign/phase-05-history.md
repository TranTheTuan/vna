# Phase 05 — Message History (Paginated List)

**Priority:** P2 | **Status:** Pending | **Effort:** 45m

## Overview

Complete the `GET /api/v1/messages` endpoint with cursor-based pagination.

## Context Links

- Brainstorm: `plans/reports/brainstorm-260403-0933-api-backend-redesign.md`
- Depends on Phase 04 (repository and service stubs exist)

## API Contract

```
GET /api/v1/messages?limit=20&cursor=<uuid>
  Header:  Authorization: Bearer <access_token>
  Success: 200 {
    "data": [
      { "id": "...", "question": "...", "answer": "...", "created_at": "..." },
      ...
    ],
    "next_cursor": "uuid-of-last-item-or-null"
  }
  Errors:
    400 invalid limit (non-integer or >100)
    401 unauthorized
```

**Cursor strategy:** UUID-based keyset pagination using `created_at` + `id` for stable ordering.
- No cursor → fetch most recent `limit` messages
- With cursor → fetch `limit` messages older than the message with that `id`

## Pagination SQL

```sql
-- No cursor (first page):
SELECT id, user_id, question, answer, created_at
FROM messages
WHERE user_id = $1
ORDER BY created_at DESC, id DESC
LIMIT $2

-- With cursor (subsequent pages):
-- First resolve cursor's created_at:
SELECT created_at FROM messages WHERE id = $1

-- Then keyset:
SELECT id, user_id, question, answer, created_at
FROM messages
WHERE user_id = $1
  AND (created_at, id) < ($cursor_created_at, $cursor_id)
ORDER BY created_at DESC, id DESC
LIMIT $2
```

> Alternatively (simpler, KISS): single query using `id < $cursor_id` UUID comparison — valid only if UUIDs are time-ordered (gen_random_uuid() v4 are NOT). Use `created_at`-based keyset above.

**`next_cursor`:** ID of the last item returned. If fewer than `limit` items returned → `next_cursor = null` (end of list).

## Files to COMPLETE (stubs from Phase 04)

### `internal/repository/message.go` — complete `ListByUser`

```go
func (r *messageRepository) ListByUser(
    ctx context.Context,
    userID string,
    limit int,
    cursor string,   // empty string = no cursor
) ([]*domain.Message, string, error)
```

Returns: `(messages, nextCursor, error)`
- `nextCursor` = last item's ID if `len(messages) == limit`, else `""`

### `internal/service/message.go` — complete `List`

```go
func (s *messageService) List(
    ctx context.Context,
    userID string,
    limit int,
    cursor string,
) ([]*domain.Message, string, error)
```

Validate: `limit` must be 1–100, default 20 if 0.
Delegate to `repo.ListByUser`.

### `internal/handler/http/message.go` — complete `List`

```go
func (h *MessageHandler) List(c echo.Context) error
```

1. Read `limit` query param (default 20, max 100)
2. Read `cursor` query param (empty = first page)
3. Get `user_id` from JWT context
4. Call `svc.List`
5. Return `ListResponse{ Data: [...], NextCursor: "..." }`

Add to `internal/dto/message_dto.go`:
```go
type ListResponse struct {
    Data       []*MessageResponse `json:"data"`
    NextCursor string             `json:"next_cursor"` // empty string → omitempty or null
}
```

## Todo List

- [ ] Complete `ListByUser` in `internal/repository/message.go`
- [ ] Complete `List` in `internal/service/message.go`
- [ ] Complete `List` handler in `internal/handler/http/message.go`
- [ ] Add `ListResponse` to `internal/dto/message_dto.go`
- [ ] `go build ./...` compiles cleanly

## Success Criteria

- First page returns ≤ `limit` items ordered newest-first
- `next_cursor` present when more pages exist, absent/null on last page
- Cursor from page N correctly fetches page N+1
- Results are isolated per user — no cross-user data leakage
- `limit > 100` returns 400

## Security Considerations

- `user_id` always taken from JWT context, never from query params
- Cursor (message ID) validated as UUID format to prevent SQL injection
