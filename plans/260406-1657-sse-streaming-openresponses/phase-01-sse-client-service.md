# Phase 01: SSE Client in Service Layer

**Priority:** High  
**Status:** ⬜ Pending  
**File:** `internal/service/message.go`

## Context

- Doc: https://docs.openclaw.ai/gateway/openresponses-http-api#streaming-sse
- Current: `callOpenResponses` sets `Stream: false`, does `io.ReadAll` → blocks until full response
- Target: `Stream: true`, reads SSE event stream, accumulates `response.output_text.delta` chunks

## SSE Protocol (from docs)

Request: `POST /v1/responses` with `stream: true`  
Response: `Content-Type: text/event-stream`

Each event line format:
```
event: <type>
data: <json>
```

Stream ends with: `data: [DONE]`

**Relevant event types:**
- `response.output_text.delta` — incremental text chunk → **extract and accumulate**
- `response.output_text.done` — full text for one output item (can use as fallback)
- `response.completed` — stream finished successfully
- `response.failed` — error during streaming

**Delta event data shape** (inferred from OpenAI Responses API compatibility):
```json
{
  "type": "response.output_text.delta",
  "delta": "chunk of text"
}
```

**Done event data shape:**
```json
{
  "type": "response.output_text.done",
  "text": "full accumulated text"
}
```

## Changes to `internal/service/message.go`

### 1. New/updated structs

```go
// openResponsesRequest — change Stream field default intent
// Stream: true for SSE

// sseEvent — parsed SSE event
type sseEvent struct {
    Type  string // from "event: <type>" line
    Data  string // from "data: <json>" line
}

// sseDeltaData — data for response.output_text.delta events
type sseDeltaData struct {
    Delta string `json:"delta"`
}

// sseDoneData — data for response.output_text.done events  
type sseDoneData struct {
    Text string `json:"text"`
}
```

### 2. New method: `streamOpenResponses`

```go
// streamOpenResponses calls POST /v1/responses with stream:true,
// reads SSE events, calls onDelta for each text chunk,
// and returns the fully-accumulated answer string.
func (s *messageService) streamOpenResponses(
    ctx context.Context,
    userID, question string,
    onDelta func(chunk string), // called for each delta, nil = accumulate only
) (string, error)
```

**Implementation steps:**
1. Build request with `Stream: true` (same auth headers as current)
2. Set `s.httpClient` — but **remove** the 30s `Timeout` on the client itself; instead rely on `ctx` deadline. The old 30s timeout kills long streams. Use `http.Client{Timeout: 0}` and pass a context with a longer deadline from the caller (handler sets it).
3. `resp.StatusCode != 200` → return `ErrUpstreamFailed`
4. Parse SSE line-by-line using `bufio.Scanner`:
   - Lines starting with `event:` → set current event type
   - Lines starting with `data:` → set current data
   - Blank line → dispatch event, reset state
   - `data: [DONE]` → break loop
5. On `response.output_text.delta`: unmarshal `sseDeltaData`, call `onDelta(delta)`, append to `var sb strings.Builder`
6. On `response.failed`: return `ErrUpstreamFailed`
7. Return `sb.String()` after stream ends

### 3. Update `MessageService` interface

Add a streaming variant:
```go
type MessageService interface {
    Send(ctx context.Context, userID, question string) (*domain.Message, error)
    SendStream(ctx context.Context, userID, question string, onDelta func(chunk string)) (*domain.Message, error)
    List(ctx context.Context, userID string, limit int, cursor string) ([]*domain.Message, string, error)
}
```

`SendStream`:
1. Calls `streamOpenResponses(ctx, userID, question, onDelta)`
2. After stream completes, saves full answer to DB via `s.repo.Save`
3. Returns saved `*domain.Message`

### 4. HTTP client timeout adjustment

The current `NewMessageService` sets `Timeout: 30 * time.Second`. This is too short for streaming. Change to:
```go
httpClient: &http.Client{
    Timeout: 0, // no global timeout; context controls cancellation
},
```

The handler will pass a context with appropriate deadline (e.g., 5 minutes).

## Implementation Steps

1. Add `sseEvent`, `sseDeltaData`, `sseDoneData` structs
2. Implement `streamOpenResponses` private method with `bufio.Scanner` SSE parser
3. Add `SendStream` to `MessageService` interface
4. Implement `SendStream` on `messageService`
5. Change `httpClient.Timeout` to `0`
6. Keep existing `Send` + `callOpenResponses` intact (backward compat)

## Todo

- [ ] Add SSE structs
- [ ] Implement `streamOpenResponses` with line-by-line SSE parsing
- [ ] Add `SendStream` to interface and implement
- [ ] Adjust `httpClient` timeout to 0
- [ ] Compile check: `go build ./...`

## Success Criteria

- `streamOpenResponses` correctly accumulates text from `response.output_text.delta` events
- Stream ends cleanly on `data: [DONE]`
- `response.failed` events surface as `ErrUpstreamFailed`
- No compile errors

## Risk

- **SSE delta JSON shape**: docs don't show exact delta payload. Assume `{"delta": "..."}` (standard OpenAI Responses API shape). Add logging for unknown event types to aid debugging.
- **bufio.Scanner default buffer**: 64KB per line. SSE data lines could be large for complex events — use `scanner.Buffer` with a larger buffer (e.g., 1MB).
