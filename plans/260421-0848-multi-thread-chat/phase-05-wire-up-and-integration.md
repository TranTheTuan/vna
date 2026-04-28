# Phase 05 — Wire-up & Integration

**Context:** `plans/reports/brainstorm-260421-0848-multi-thread-chat.md`
**Priority:** P1 | **Status:** Pending | **Effort:** 0.5h
**Blocked by:** Phase 04

## Overview

Wire `ThreadRepository`, `ThreadService`, and `ThreadHandler` into `main.go`. Update `MessageService` constructor to receive `ThreadRepository`.

## Related Code Files

| Action | File |
|--------|------|
| Modify | `cmd/api/main.go` |
| Modify | `internal/delivery/http/message_routes.go` |

## Implementation Steps

### 1. `cmd/api/main.go`

Add after existing repo/service/handler construction:

```go
// Repositories
userRepo    := repository.NewUserRepository(pool)
messageRepo := repository.NewMessageRepository(pool)
threadRepo  := repository.NewThreadRepository(pool)   // ← new

// Services
userSvc    := service.NewUserService(cfg, userRepo, logger)
messageSvc := service.NewMessageService(cfg, messageRepo, threadRepo, logger)  // ← add threadRepo
threadSvc  := service.NewThreadService(threadRepo)    // ← new

// Handlers
authHandler    := http_handler.NewAuthHandler(userSvc)
messageHandler := http_handler.NewMessageHandler(messageSvc)
threadHandler  := http_handler.NewThreadHandler(threadSvc)   // ← new

// Routes
http_delivery.RegisterAuthRoutes(apiGroup, authHandler, cfg)
http_delivery.RegisterMessageRoutes(apiGroup, messageHandler, cfg)
http_delivery.RegisterThreadRoutes(apiGroup, threadHandler, cfg)  // ← new
```

### 2. `NewMessageService` signature update

```go
func NewMessageService(
    cfg *configs.Config,
    repo repository.MessageRepository,
    threadRepo repository.ThreadRepository,  // ← new param
    logger *slog.Logger,
) MessageService {
    return &messageService{
        cfg:        cfg,
        repo:       repo,
        threadRepo: threadRepo,
        httpClient: &http.Client{},
        logger:     logger,
    }
}
```

### 3. Check `internal/delivery/http/message_routes.go`

Ensure `GET /messages` route references `h.ListByThread` (renamed in Phase 04). No other changes needed.

## Todo

- [x] Update `NewMessageService` constructor to accept `threadRepo`
- [x] Add `threadRepo repository.ThreadRepository` field to `messageService` struct
- [x] Add `threadRepo`, `threadSvc`, `threadHandler` to `cmd/api/main.go`
- [x] Register `RegisterThreadRoutes` in `cmd/api/main.go`
- [x] Verify `GET /messages` route uses `ListByThread`
- [x] Run `go build ./...` — must pass with zero errors

## Success Criteria

- `go build ./...` clean
- Server starts without panic
- `GET /api/v1/threads` reachable (returns empty array for new user)
- `POST /api/v1/messages/stream` with empty `thread_id` creates thread and emits `event: metadata`
