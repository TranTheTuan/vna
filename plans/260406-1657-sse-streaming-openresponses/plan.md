# Plan: SSE Streaming for OpenResponses API

**Created:** 2026-04-06  
**Status:** ✅ Complete  
**Branch:** main

## Problem

`callOpenResponses` in `internal/service/message.go` uses `Stream: false` (blocking request). Long AI responses cause:
- Client 504 Gateway Timeout
- Poor UX (no feedback until full response arrives)

## Solution

Switch `POST /api/v1/messages` to SSE streaming:
1. Set `stream: true` → OpenResponses sends `text/event-stream`
2. Service reads `response.output_text.delta` events, accumulates full text
3. Handler proxies SSE chunks to client in real-time (before DB save)
4. DB save happens after stream completes

## Phases

| # | Phase | Status |
|---|-------|--------|
| 1 | [SSE client in service layer](phase-01-sse-client-service.md) | ✅ Complete |
| 2 | [SSE handler & new endpoint](phase-02-sse-handler-endpoint.md) | ✅ Complete |
| 3 | [Router & Swagger update](phase-03-router-swagger.md) | ✅ Complete |

## Key Files

- `internal/service/message.go` — core change (SSE client, stream accumulator)
- `internal/handler/http/message.go` — new `SendStream` handler
- `internal/delivery/http/router.go` — register new route
- `internal/dto/message_dto.go` — SSE chunk DTO

## Architecture Decision

Two endpoints coexist:
- `POST /api/v1/messages` — existing (keep for backward compat)
- `POST /api/v1/messages/stream` — **new** SSE endpoint

This avoids breaking existing clients while enabling streaming for new ones.
