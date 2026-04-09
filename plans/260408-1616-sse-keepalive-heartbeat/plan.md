---
title: "SSE Keepalive Heartbeat Goroutine"
description: "Add heartbeat goroutine to SendStream handler to prevent Cloudflare 524 timeout during long AI thinking gaps"
status: completed
priority: P1
effort: 1h
issue:
branch: feat/send-sse
tags: [backend, api, feature, critical]
created: 2026-04-08
---

# SSE Keepalive Heartbeat Goroutine

## Overview

The `POST /api/v1/messages/stream` SSE endpoint can trigger Cloudflare 524 errors when the upstream AI model has long thinking gaps (>100s) between delta events. Add a heartbeat goroutine inside `SendStream` that emits `event: ping` every 15s when no delta was recently flushed.

## Phases

| # | Phase | Status | Effort | Link |
|---|-------|--------|--------|------|
| 1 | Implement heartbeat goroutine | ✅ Complete | 45m | [phase-01](./phase-01-heartbeat-goroutine.md) |
| 2 | Test & verify | ✅ Complete | 15m | [phase-02-test.md](./phase-02-test.md) |

## Dependencies

- Only `internal/handler/http/message.go` modified
- No service, DTO, or client changes
