# Phase 01: Implement Heartbeat Goroutine

**Priority:** P1  
**Status:** ⬜ Pending  
**Effort:** 45m

## Context

- Problem: CF 524 fires when upstream AI has >100s thinking gap (no SSE events emitted)
- Solution: Goroutine fires `event: ping\ndata: {}\n\n` every 15s if no delta flushed recently
- Brainstorm report: `plans/reports/brainstorm-260408-1616-sse-keepalive-heartbeat.md`
- Handler file: `internal/handler/http/message.go` (currently 220 lines)

## Key Insights

- CF 524 triggers when **server-side** connection is silent >100s — client-side EventSource reconnect doesn't help
- Heartbeat must be serialized with `onDelta` writes to prevent interleaved SSE frames
- `atomic.Int64` for `lastFlush` timestamp — avoids mutex for the timestamp check, only lock the writer
- `done` channel closed via `defer close(done)` ensures goroutine never leaks regardless of error paths
- Standard SSE clients (`EventSource`, fetch) silently ignore `event: ping` unless explicitly subscribed
- 15s interval → 6x safety margin under CF's 100s limit; check threshold 14s to account for ticker jitter

## Requirements

- Emit `event: ping\ndata: {}\n\n` when no flush in >14s
- Goroutine must not outlive `SendStream` handler call
- No data race between heartbeat flush and `onDelta` flush
- No changes to service layer, DTOs, or client

## Architecture

```
SendStream handler
│
├── sync.Mutex (mu)           ← serializes all writer access
├── atomic.Int64 (lastFlush)  ← timestamp of last successful flush
├── chan struct{} (done)       ← signals goroutine shutdown
│
├── goroutine: heartbeat
│   └── ticker 15s → if time.Since(lastFlush) > 14s → mu.Lock → ping → Unlock
│
└── onDelta callback
    └── mu.Lock → writeSseEvent(delta) → Unlock → lastFlush.Store(now)
```

## Related Code Files

| File | Action | Description |
|------|--------|-------------|
| `internal/handler/http/message.go` | Modify | Add heartbeat goroutine to `SendStream` only |

## Implementation Steps

1. Add imports `sync` and `sync/atomic` to `internal/handler/http/message.go`
   - `sync` already imported? check — if not, add it
   - `sync/atomic` — add to imports

2. In `SendStream`, after asserting `flusher` support and writing SSE headers, add:

```go
// Track last flush time for keepalive heartbeat.
var mu sync.Mutex
var lastFlush atomic.Int64
lastFlush.Store(time.Now().UnixNano())

// Heartbeat goroutine: emits "ping" every 15s if no delta was flushed recently.
// Prevents Cloudflare 524 timeout during long upstream thinking gaps.
done := make(chan struct{})
defer close(done)
go func() {
    ticker := time.NewTicker(15 * time.Second)
    defer ticker.Stop()
    for {
        select {
        case <-ticker.C:
            if time.Since(time.Unix(0, lastFlush.Load())) > 14*time.Second {
                mu.Lock()
                writeSseEvent(w, flusher, "ping", "{}")
                mu.Unlock()
            }
        case <-done:
            return
        }
    }
}()
```

3. Wrap `onDelta` callback to use mutex and update `lastFlush`:

```go
onDelta := func(chunk string) {
    data, _ := json.Marshal(dto.StreamDeltaEvent{Delta: chunk})
    mu.Lock()
    writeSseEvent(w, flusher, "delta", string(data))
    mu.Unlock()
    lastFlush.Store(time.Now().UnixNano())
}
```

4. Wrap final `writeSseEvent` calls (error + done events) with mutex:

```go
// error path:
mu.Lock()
writeSseEvent(w, flusher, "error", string(errData))
mu.Unlock()

// done path:
mu.Lock()
writeSseEvent(w, flusher, "done", string(doneData))
mu.Unlock()
```

5. Verify file stays under 200 lines after changes

6. Run `go build ./...` to confirm no compile errors

## Todo

- [x] Add `sync/atomic` import
- [x] Add `sync.Mutex` + `atomic.Int64` vars after header write
- [x] Add heartbeat goroutine with `done` channel
- [x] Update `onDelta` to use mutex + update `lastFlush`
- [x] Wrap error/done `writeSseEvent` calls with mutex
- [x] Verify line count <200
- [x] Run `go build ./...`

## Success Criteria

- `go build ./...` passes with no errors
- `go vet ./...` passes with no race warnings
- `SendStream` handler ≤200 lines
- Heartbeat goroutine never leaks (verified by `defer close(done)`)
- No data race between heartbeat and `onDelta` (mutex serializes both)

## Risk Assessment

| Risk | Likelihood | Mitigation |
|------|-----------|------------|
| Goroutine leak if panic in handler | Low | `defer close(done)` runs even on panic recovery |
| `writeSseEvent` after client disconnect | Low | Go ignores writes to closed connections; no panic |
| Ticker fires after stream ends | None | `done` channel unblocks select immediately |

## Security Considerations

- `ping` event emits only `{}` — no user data leaked
- No auth/session info in keepalive payload
