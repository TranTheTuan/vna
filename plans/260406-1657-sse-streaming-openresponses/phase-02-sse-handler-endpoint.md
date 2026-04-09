# Phase 02: SSE Handler & New Endpoint

**Priority:** High  
**Status:** ⬜ Pending  
**Depends on:** Phase 01  
**Files:** `internal/handler/http/message.go`, `internal/dto/message_dto.go`

## Context

- Phase 01 adds `SendStream` to `MessageService` with an `onDelta func(chunk string)` callback
- This phase wires that callback to flush SSE chunks directly to the HTTP response
- Echo supports SSE via `c.Response().Writer` + manual flushing

## SSE Response Format (client-facing)

The handler will emit standard SSE to the browser/client:

```
Content-Type: text/event-stream
Cache-Control: no-cache
Connection: keep-alive

event: delta
data: {"delta":"Hello"}

event: delta
data: {"delta":" world"}

event: done
data: {"id":"uuid","question":"...","answer":"full text","created_at":"..."}

```

- `event: delta` — incremental chunk while streaming
- `event: done` — final event with complete saved message (mirrors existing `MessageResponse`)
- `event: error` — on upstream failure

## Changes to `internal/dto/message_dto.go`

Add one struct:

```go
// StreamDeltaEvent is sent as SSE "delta" events during streaming.
type StreamDeltaEvent struct {
    Delta string `json:"delta"`
}
```

No changes to existing DTOs.

## Changes to `internal/handler/http/message.go`

### New handler: `SendStream`

```go
// SendStream handles POST /api/v1/messages/stream.
// Returns a text/event-stream SSE response.
// Streams delta chunks while the AI responds, then emits a final "done"
// event with the saved message once the stream completes.
//
// @Summary      Stream a chat message
// @Description  Sends a message to the AI and streams the response via SSE.
// @Tags         messages
// @Accept       json
// @Produce      text/event-stream
// @Security     BearerAuth
// @Param        body  body  dto.SendMessageRequest  true  "Message request"
// @Success      200   {string} string "SSE stream"
// @Failure      400   {object} map[string]string
// @Failure      401   {object} map[string]string
// @Failure      502   {object} map[string]string
// @Router       /api/v1/messages/stream [post]
func (h *MessageHandler) SendStream(c echo.Context) error
```

**Implementation:**

1. Bind & validate `dto.SendMessageRequest` (same as `Send`)
2. Extract `userID` from context
3. Set SSE headers:
   ```go
   c.Response().Header().Set("Content-Type", "text/event-stream")
   c.Response().Header().Set("Cache-Control", "no-cache")
   c.Response().Header().Set("Connection", "keep-alive")
   c.Response().WriteHeader(http.StatusOK)
   ```
4. Get `http.Flusher` from `c.Response().Writer`:
   ```go
   flusher, ok := c.Response().Writer.(http.Flusher)
   if !ok {
       return echo.NewHTTPError(http.StatusInternalServerError, "streaming not supported")
   }
   ```
5. Create a 5-minute context:
   ```go
   ctx, cancel := context.WithTimeout(c.Request().Context(), 5*time.Minute)
   defer cancel()
   ```
6. Build `onDelta` callback — marshals `StreamDeltaEvent`, writes SSE line, flushes:
   ```go
   onDelta := func(chunk string) {
       data, _ := json.Marshal(dto.StreamDeltaEvent{Delta: chunk})
       fmt.Fprintf(c.Response().Writer, "event: delta\ndata: %s\n\n", data)
       flusher.Flush()
   }
   ```
7. Call `h.svc.SendStream(ctx, userID, req.Message, onDelta)`
8. On success: marshal `MessageResponse`, emit `event: done`, flush
9. On error: emit `event: error` with message text, flush (do NOT return an HTTP error — headers already sent)

**Error handling after headers are sent:**

Once `WriteHeader(200)` is called, HTTP error codes can't be changed. Use SSE error events:
```
event: error
data: {"message":"upstream API error"}

```

### Helper: `writeSseEvent`

Extract a small private helper to avoid repetition:
```go
func writeSseEvent(w io.Writer, f http.Flusher, event, data string) {
    fmt.Fprintf(w, "event: %s\ndata: %s\n\n", event, data)
    f.Flush()
}
```

## Implementation Steps

1. Add `StreamDeltaEvent` to `internal/dto/message_dto.go`
2. Add `writeSseEvent` helper in `internal/handler/http/message.go`
3. Implement `SendStream` handler
4. Compile check: `go build ./...`

## Todo

- [ ] Add `StreamDeltaEvent` to dto
- [ ] Add `writeSseEvent` private helper
- [ ] Implement `SendStream` handler with SSE headers + flush loop
- [ ] Compile check

## Success Criteria

- Handler sets correct SSE headers before writing any body
- Each `onDelta` call emits one `event: delta` line and flushes immediately
- Final `event: done` contains full `MessageResponse` JSON
- Errors after stream start emit `event: error` (no panic, no broken pipe crash)
- No compile errors

## Risk

- **Echo response writer flushing**: Echo wraps `http.ResponseWriter`; `c.Response().Writer` is the underlying writer. Confirm `http.Flusher` assertion works — it does in standard `net/http` servers.
- **Concurrent writes**: `onDelta` is called synchronously inside `streamOpenResponses` (single goroutine), so no mutex needed.
