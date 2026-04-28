# Phase 03 — Service Layer

**Context:** `plans/reports/brainstorm-260421-0848-multi-thread-chat.md`
**Priority:** P1 | **Status:** Pending | **Effort:** 1.5h
**Blocked by:** Phase 02

## Overview

Add `ThreadService` for thread CRUD. Update `MessageService` to accept `threadID` on send/stream/list, handle implicit thread creation, emit `metadata` SSE event first.

## Requirements

- `ThreadService`: `Create(userID)`, `ListByUser(userID)` 
- `MessageService.Send/SendStream`: accept `threadID string` — empty = create new thread
- `MessageService.List` replaced by `ListByThread(userID, threadID, limit, cursor)`
- OpenResponses `user` field = `threadID` (not `userID`)
- New sentinel error: `ErrThreadNotFound`

## Related Code Files

| Action | File |
|--------|------|
| Create | `internal/service/thread.go` |
| Modify | `internal/service/message.go` |

## Implementation Steps

### 1. `internal/service/thread.go`

```go
package service

import (
    "context"
    "database/sql"
    "errors"
    "fmt"

    "github.com/TranTheTuan/vna/internal/domain"
    "github.com/TranTheTuan/vna/internal/repository"
)

var ErrThreadNotFound = errors.New("service.thread: thread not found")

type ThreadService interface {
    Create(ctx context.Context, userID string) (*domain.Thread, error)
    ListByUser(ctx context.Context, userID string) ([]*domain.Thread, error)
}

type threadService struct {
    repo repository.ThreadRepository
}

func NewThreadService(repo repository.ThreadRepository) ThreadService {
    return &threadService{repo: repo}
}

func (s *threadService) Create(ctx context.Context, userID string) (*domain.Thread, error) {
    return s.repo.Create(ctx, userID)
}

func (s *threadService) ListByUser(ctx context.Context, userID string) ([]*domain.Thread, error) {
    return s.repo.ListByUser(ctx, userID)
}
```

### 2. `internal/service/message.go` — interface changes

Update `MessageService` interface:

```go
type MessageService interface {
    Send(ctx context.Context, userID, threadID, question string) (*domain.Message, error)
    SendStream(ctx context.Context, userID, threadID, question string, onMeta func(threadID string), onDelta func(chunk string)) (*domain.Message, error)
    ListByThread(ctx context.Context, userID, threadID string, limit int, cursor string) ([]*domain.Message, string, error)
}
```

Key changes:
- `Send`/`SendStream`: add `threadID string` param
- `SendStream`: add `onMeta func(threadID string)` callback — handler calls this to emit `event: metadata`
- `List` → `ListByThread` (scoped to thread, validates ownership)

### 3. Thread resolution helper (inside `messageService`)

```go
// resolveThread returns existing threadID or creates a new thread.
// Validates ownership when threadID is non-empty.
func (s *messageService) resolveThread(ctx context.Context, userID, threadID string) (string, error) {
    if threadID == "" {
        t, err := s.threadRepo.Create(ctx, userID)
        if err != nil {
            return "", fmt.Errorf("service.message: create thread: %w", err)
        }
        return t.ID, nil
    }
    // Validate ownership
    _, err := s.threadRepo.GetByIDAndUser(ctx, threadID, userID)
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return "", ErrThreadNotFound
        }
        return "", fmt.Errorf("service.message: validate thread: %w", err)
    }
    return threadID, nil
}
```

### 4. `SendStream` updated flow

```go
func (s *messageService) SendStream(ctx context.Context, userID, threadID, question string,
    onMeta func(string), onDelta func(string)) (*domain.Message, error) {
    if question == "" {
        return nil, ErrEmptyMessage
    }
    resolvedID, err := s.resolveThread(ctx, userID, threadID)
    if err != nil {
        return nil, err
    }
    // Notify handler of thread ID before streaming starts
    if onMeta != nil {
        onMeta(resolvedID)
    }
    // Pass thread ID as OpenResponses 'user' param
    answer, err := s.streamOpenResponses(ctx, resolvedID, question, onDelta)
    if err != nil { ... }
    msg, err := s.repo.Save(ctx, &domain.Message{
        UserID:   userID,
        ThreadID: resolvedID,
        Question: question,
        Answer:   answer,
    })
    ...
}
```

### 5. `openResponsesRequest` — `User` field stays, value changes

The struct is unchanged. The call site changes from `User: userID` to `User: threadID`.

### 6. `ListByThread`

```go
func (s *messageService) ListByThread(ctx context.Context, userID, threadID string, limit int, cursor string) ([]*domain.Message, string, error) {
    if limit == 0 { limit = 20 }
    if limit < 1 || limit > 100 { return nil, "", ErrInvalidLimit }
    // Validate thread ownership
    if _, err := s.threadRepo.GetByIDAndUser(ctx, threadID, userID); err != nil {
        if errors.Is(err, sql.ErrNoRows) { return nil, "", ErrThreadNotFound }
        return nil, "", err
    }
    return s.repo.ListByThread(ctx, threadID, limit, cursor)
}
```

## Todo

- [x] Create `internal/service/thread.go`
- [x] Add `threadRepo repository.ThreadRepository` field to `messageService` struct
- [x] Add `resolveThread` helper to `messageService`
- [x] Update `Send` signature and logic
- [x] Update `SendStream` signature — add `onMeta` callback, call `onMeta(resolvedID)` before streaming
- [x] Replace `callOpenResponses`/`streamOpenResponses` `user` arg: `userID` → `threadID`
- [x] Replace `List` → `ListByThread` in interface + implementation
- [x] Add `ErrThreadNotFound` sentinel

## Success Criteria

- `go build ./...` passes
- `messageService` satisfies updated `MessageService` interface
- `threadService` satisfies `ThreadService` interface

## Security Considerations

- `GetByIDAndUser` always validates user owns the thread before any operation — prevents cross-user thread access
- Thread creation uses `userID` from JWT context, never from request body
