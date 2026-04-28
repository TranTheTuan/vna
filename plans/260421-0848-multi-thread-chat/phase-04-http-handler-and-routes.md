# Phase 04 — HTTP Handler & Routes

**Context:** `plans/reports/brainstorm-260421-0848-multi-thread-chat.md`
**Priority:** P1 | **Status:** Pending | **Effort:** 1.5h
**Blocked by:** Phase 03

## Overview

Add `ThreadHandler` + routes. Update `MessageHandler` to pass `thread_id` through and emit `event: metadata` SSE. Update DTOs.

## Related Code Files

| Action | File |
|--------|------|
| Modify | `internal/dto/message_dto.go` |
| Create | `internal/dto/thread_dto.go` |
| Modify | `internal/handler/http/message.go` |
| Create | `internal/handler/http/thread.go` |
| Create | `internal/delivery/http/thread_routes.go` |

## Implementation Steps

### 1. `internal/dto/message_dto.go` — extend existing DTOs

```go
// SendMessageRequest — add optional ThreadID
type SendMessageRequest struct {
    Message  string `json:"message"`
    ThreadID string `json:"thread_id"` // empty = create new thread
}

// MessageResponse — add ThreadID
type MessageResponse struct {
    ID        string    `json:"id"`
    ThreadID  string    `json:"thread_id"`
    Question  string    `json:"question"`
    Answer    string    `json:"answer"`
    CreatedAt time.Time `json:"created_at"`
}

// StreamMetaEvent — emitted as first SSE event
type StreamMetaEvent struct {
    ThreadID string `json:"thread_id"`
}

// StreamDeltaEvent unchanged
// ListResponse unchanged
```

### 2. `internal/dto/thread_dto.go`

```go
package dto

import "time"

type ThreadResponse struct {
    ID        string    `json:"id"`
    Name      string    `json:"name"`
    CreatedAt time.Time `json:"created_at"`
}

type ListThreadsResponse struct {
    Data []*ThreadResponse `json:"data"`
}
```

### 3. `internal/handler/http/message.go` — update `SendStream` and `List`

**`SendStream` changes:**
```go
func (h *MessageHandler) SendStream(c echo.Context) error {
    var req dto.SendMessageRequest
    if err := c.Bind(&req); err != nil { ... }
    if req.Message == "" { return echo.NewHTTPError(400, "message is required") }

    userID := c.Get("user_id").(string)
    flusher, ok := c.Response().Writer.(http.Flusher)
    if !ok { ... }

    c.Response().Header().Set("Content-Type", "text/event-stream")
    c.Response().Header().Set("Cache-Control", "no-cache")
    c.Response().Header().Set("Connection", "keep-alive")
    c.Response().WriteHeader(http.StatusOK)

    ctx, cancel := context.WithTimeout(c.Request().Context(), 5*time.Minute)
    defer cancel()

    w := c.Response().Writer

    // onMeta emits thread_id as FIRST event, before any delta
    onMeta := func(threadID string) {
        data, _ := json.Marshal(dto.StreamMetaEvent{ThreadID: threadID})
        writeSseEvent(w, flusher, "metadata", string(data))
    }
    onDelta := func(chunk string) {
        data, _ := json.Marshal(dto.StreamDeltaEvent{Delta: chunk})
        writeSseEvent(w, flusher, "delta", string(data))
    }

    msg, err := h.svc.SendStream(ctx, userID, req.ThreadID, req.Message, onMeta, onDelta)
    if err != nil {
        errData, _ := json.Marshal(map[string]string{"message": sseErrorMessage(err)})
        writeSseEvent(w, flusher, "error", string(errData))
        return nil
    }

    doneData, _ := json.Marshal(dto.MessageResponse{
        ID:        msg.ID,
        ThreadID:  msg.ThreadID,
        Question:  msg.Question,
        Answer:    msg.Answer,
        CreatedAt: msg.CreatedAt,
    })
    writeSseEvent(w, flusher, "done", string(doneData))
    return nil
}
```

**`List` changes:** rename to `ListByThread`, require `?thread_id` query param:
```go
func (h *MessageHandler) ListByThread(c echo.Context) error {
    userID := c.Get("user_id").(string)
    threadID := c.QueryParam("thread_id")
    if threadID == "" {
        return echo.NewHTTPError(http.StatusBadRequest, "thread_id is required")
    }
    // limit/cursor parsing unchanged ...
    msgs, nextCursor, err := h.svc.ListByThread(c.Request().Context(), userID, threadID, limit, cursor)
    if err != nil {
        switch {
        case errors.Is(err, service.ErrThreadNotFound):
            return echo.NewHTTPError(http.StatusNotFound, "thread not found")
        case errors.Is(err, service.ErrInvalidLimit):
            return echo.NewHTTPError(http.StatusBadRequest, "limit must be between 1 and 100")
        default:
            return echo.NewHTTPError(http.StatusInternalServerError, "failed to list messages")
        }
    }
    // build response same as before, include ThreadID in each MessageResponse
    ...
}
```

Also update `Send` to accept and pass `req.ThreadID`. Add `ErrThreadNotFound` case to `sseErrorMessage`.

### 4. `internal/handler/http/thread.go`

```go
package http

import (
    "net/http"
    "github.com/labstack/echo/v4"
    "github.com/TranTheTuan/vna/internal/dto"
    "github.com/TranTheTuan/vna/internal/service"
)

type ThreadHandler struct {
    svc service.ThreadService
}

func NewThreadHandler(svc service.ThreadService) *ThreadHandler {
    return &ThreadHandler{svc: svc}
}

// List handles GET /api/v1/threads
func (h *ThreadHandler) List(c echo.Context) error {
    userID := c.Get("user_id").(string)
    threads, err := h.svc.ListByUser(c.Request().Context(), userID)
    if err != nil {
        return echo.NewHTTPError(http.StatusInternalServerError, "failed to list threads")
    }
    data := make([]*dto.ThreadResponse, len(threads))
    for i, t := range threads {
        data[i] = &dto.ThreadResponse{ID: t.ID, Name: t.Name, CreatedAt: t.CreatedAt}
    }
    return c.JSON(http.StatusOK, dto.ListThreadsResponse{Data: data})
}
```

### 5. `internal/delivery/http/thread_routes.go`

```go
package http

import (
    "github.com/labstack/echo/v4"
    handler "github.com/TranTheTuan/vna/internal/handler/http"
    "github.com/TranTheTuan/vna/configs"
    "github.com/TranTheTuan/vna/pkg/jwt_util"
)

func RegisterThreadRoutes(g *echo.Group, h *handler.ThreadHandler, cfg *configs.Config) {
    threads := g.Group("/threads", jwt_util.JWTMiddleware(cfg))
    threads.GET("", h.List)
}
```

Update `internal/delivery/http/message_routes.go`:
- Route `GET /messages` calls `h.ListByThread` (rename handler method reference)

## Todo

- [x] Add `ThreadID` to `SendMessageRequest` and `MessageResponse` in `message_dto.go`
- [x] Add `StreamMetaEvent` to `message_dto.go`
- [x] Create `internal/dto/thread_dto.go`
- [x] Update `SendStream` in `message.go` — add `onMeta` callback, pass `req.ThreadID`
- [x] Rename `List` → `ListByThread` in `message.go`, add `thread_id` query param validation
- [x] Update `Send` in `message.go` to pass `req.ThreadID`
- [x] Add `ErrThreadNotFound` case to `sseErrorMessage`
- [x] Create `internal/handler/http/thread.go`
- [x] Create `internal/delivery/http/thread_routes.go`
- [x] Update message routes to use `ListByThread`

## Success Criteria

- `go build ./...` passes
- `GET /api/v1/threads` returns 200 with thread list
- `POST /api/v1/messages/stream` emits `event: metadata` as first SSE event
- `GET /api/v1/messages?thread_id=xxx` scopes results to thread

## Security Considerations

- All thread routes behind JWT middleware
- `thread_id` from request body/query is always validated against `user_id` from JWT — never trust client-supplied thread ownership
