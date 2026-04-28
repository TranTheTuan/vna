# Phase 02 — Domain & Repository Layer

**Context:** `plans/reports/brainstorm-260421-0848-multi-thread-chat.md`
**Priority:** P1 | **Status:** Pending | **Effort:** 1h
**Blocked by:** Phase 01 (migration must be applied)

## Overview

Add `Thread` domain type, `ThreadRepository` interface + implementation. Update `Message` domain and `MessageRepository` to be thread-scoped.

## Requirements

- `domain.Thread` struct
- `ThreadRepository` interface: `Create`, `ListByUser`, `GetByID`
- `Message.ThreadID` field added
- `MessageRepository.Save` accepts `thread_id`
- `MessageRepository.ListByUser` replaced by `ListByThread(ctx, threadID, limit, cursor)`

## Related Code Files

| Action | File |
|--------|------|
| Create | `internal/domain/thread.go` |
| Modify | `internal/domain/message.go` |
| Create | `internal/repository/thread.go` |
| Modify | `internal/repository/message.go` |

## Implementation Steps

### 1. `internal/domain/thread.go`

```go
package domain

import "time"

// Thread represents a named chat conversation belonging to a user.
type Thread struct {
    ID        string
    UserID    string
    Name      string
    CreatedAt time.Time
}
```

### 2. `internal/domain/message.go` — add `ThreadID` field

```go
type Message struct {
    ID        string
    UserID    string
    ThreadID  string    // UUID of owning thread
    Question  string
    Answer    string
    CreatedAt time.Time
}
```

### 3. `internal/repository/thread.go`

```go
package repository

import (
    "context"
    "database/sql"
    "fmt"

    "github.com/TranTheTuan/vna/internal/domain"
)

type ThreadRepository interface {
    Create(ctx context.Context, userID string) (*domain.Thread, error)
    ListByUser(ctx context.Context, userID string) ([]*domain.Thread, error)
    GetByIDAndUser(ctx context.Context, threadID, userID string) (*domain.Thread, error)
}

type threadRepository struct{ db *sql.DB }

func NewThreadRepository(db *sql.DB) ThreadRepository {
    return &threadRepository{db: db}
}

func (r *threadRepository) Create(ctx context.Context, userID string) (*domain.Thread, error) {
    const q = `
        INSERT INTO threads(user_id)
        VALUES($1)
        RETURNING id, user_id, name, created_at`
    t := &domain.Thread{}
    err := r.db.QueryRowContext(ctx, q, userID).
        Scan(&t.ID, &t.UserID, &t.Name, &t.CreatedAt)
    if err != nil {
        return nil, fmt.Errorf("repository.thread: create: %w", err)
    }
    return t, nil
}

func (r *threadRepository) ListByUser(ctx context.Context, userID string) ([]*domain.Thread, error) {
    const q = `
        SELECT id, user_id, name, created_at
        FROM threads
        WHERE user_id = $1
        ORDER BY created_at DESC`
    rows, err := r.db.QueryContext(ctx, q, userID)
    if err != nil {
        return nil, fmt.Errorf("repository.thread: list: %w", err)
    }
    defer rows.Close()
    var threads []*domain.Thread
    for rows.Next() {
        t := &domain.Thread{}
        if err := rows.Scan(&t.ID, &t.UserID, &t.Name, &t.CreatedAt); err != nil {
            return nil, fmt.Errorf("repository.thread: scan: %w", err)
        }
        threads = append(threads, t)
    }
    return threads, rows.Err()
}

func (r *threadRepository) GetByIDAndUser(ctx context.Context, threadID, userID string) (*domain.Thread, error) {
    const q = `SELECT id, user_id, name, created_at FROM threads WHERE id=$1 AND user_id=$2`
    t := &domain.Thread{}
    err := r.db.QueryRowContext(ctx, q, threadID, userID).
        Scan(&t.ID, &t.UserID, &t.Name, &t.CreatedAt)
    if err != nil {
        return nil, fmt.Errorf("repository.thread: get: %w", err)
    }
    return t, nil
}
```

### 4. `internal/repository/message.go` — update `Save` and replace `ListByUser` with `ListByThread`

- `Save`: add `thread_id` to INSERT and RETURNING scan
- Remove `ListByUser`; add `ListByThread(ctx, threadID, userID, limit, cursor)`:
  - First page: `WHERE thread_id=$1` + keyset on `(created_at DESC, id DESC)`
  - Cursor page: resolve cursor `created_at`, then `WHERE thread_id=$1 AND (created_at,id) < ($2,$3)`
  - `userID` used only for cursor validation (cursor must belong to same thread)

> Note: `ListByThread` drops the `userID` filter on main rows — thread ownership is already validated at service layer before this call.

## Todo

- [x] Create `internal/domain/thread.go`
- [x] Add `ThreadID` to `internal/domain/message.go`
- [x] Create `internal/repository/thread.go`
- [x] Update `Save` in `internal/repository/message.go` (add thread_id col)
- [x] Replace `ListByUser` → `ListByThread` in `internal/repository/message.go`

## Success Criteria

- All files compile (`go build ./...`)
- `ThreadRepository` and updated `MessageRepository` interfaces are satisfied by their implementations

## Risk Assessment

- `ListByUser` removal is a breaking change at the interface level — update all call sites in Phase 03
