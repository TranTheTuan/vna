---
title: "Multi-Thread Chat"
description: "Add thread-scoped chat history — users can create multiple named threads, each with isolated AI context."
status: in-progress
priority: P1
effort: 6h
issue:
branch: feat/thread-message
tags: [feature, backend, database, api]
created: 2026-04-21
---

# Multi-Thread Chat

## Overview

Currently one flat chat history per user. This plan adds thread-based isolation:
- Users create threads implicitly (empty `thread_id` on stream call)
- Thread ID streamed back as first SSE event (`event: metadata`)
- Messages scoped to threads; OpenResponses `user` param = `thread_id`

See brainstorm report: `plans/reports/brainstorm-260421-0848-multi-thread-chat.md`

## Phases

| # | Phase | Status | Effort | Link |
|---|-------|--------|--------|------|
| 1 | Database migration | Complete | 0.5h | [phase-01](./phase-01-database-migration.md) |
| 2 | Domain & repository layer | Complete | 1h | [phase-02-domain-and-repository.md](./phase-02-domain-and-repository.md) |
| 3 | Service layer | Complete | 1.5h | [phase-03-service-layer.md](./phase-03-service-layer.md) |
| 4 | HTTP handler & routes | Complete | 1.5h | [phase-04-http-handler-and-routes.md](./phase-04-http-handler-and-routes.md) |
| 5 | Wire-up & integration | Complete | 0.5h | [phase-05-wire-up-and-integration.md](./phase-05-wire-up-and-integration.md) |
| 6 | Tests | Complete | 1h | [phase-06-tests.md](./phase-06-tests.md) |

## Dependencies

- Phases 1 → 2 → 3 → 4 → 5 → 6 (strict sequential)
- Phase 2 depends on migration (schema must exist before repo layer)
- Breaking change: `GET /api/v1/messages` now requires `?thread_id` — coordinate with frontend
