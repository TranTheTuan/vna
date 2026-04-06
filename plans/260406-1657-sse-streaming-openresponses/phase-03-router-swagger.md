# Phase 03: Router & Swagger Update

**Priority:** High  
**Status:** ⬜ Pending  
**Depends on:** Phase 02  
**Files:** `internal/delivery/http/router.go`, `internal/docs/docs.go` (auto-generated)

## Context

- Phase 02 adds `SendStream` to `MessageHandler`
- This phase registers the new route and regenerates Swagger docs

## Changes to `internal/delivery/http/router.go`

Add one line to `RegisterRoutes`:

```go
msgs.POST("/stream", mh.SendStream)
```

Full updated block:
```go
msgs := e.Group("/api/v1/messages")
msgs.Use(JWTMiddleware(cfg))
msgs.POST("", mh.Send)
msgs.GET("", mh.List)
msgs.POST("/stream", mh.SendStream)  // ← new
```

## Swagger Regeneration

The `SendStream` handler already has `@Router /api/v1/messages/stream [post]` annotations from Phase 02.

Run:
```bash
swag init -g cmd/main.go -o internal/docs
```

This regenerates `internal/docs/docs.go`, `swagger.json`, `swagger.yaml`.

> **Note:** `internal/docs/docs.go` is auto-generated — never edit manually.

## Implementation Steps

1. Add `msgs.POST("/stream", mh.SendStream)` to router
2. Run `swag init` to regenerate docs
3. Compile check: `go build ./...`
4. Manual smoke test (optional):
   ```bash
   curl -N -X POST http://localhost:8080/api/v1/messages/stream \
     -H "Authorization: Bearer <token>" \
     -H "Content-Type: application/json" \
     -d '{"message":"hello"}' --no-buffer
   ```
   Expect: SSE `event: delta` lines streaming in, then `event: done`.

## Todo

- [ ] Add `/stream` route to router
- [ ] Run `swag init -g cmd/main.go -o internal/docs`
- [ ] Compile check: `go build ./...`

## Success Criteria

- `POST /api/v1/messages/stream` route registered and reachable
- Swagger UI shows new endpoint
- No compile errors
- Existing `POST /api/v1/messages` unaffected
