---
title: "API Backend Redesign — Auth + Chat + History"
description: "Remove IoT/TCP/OAuth code; build email/password JWT auth, OpenResponses-backed chat, and paginated message history on Go/Echo/pgx."
status: pending
priority: P1
effort: 6h
issue:
branch: main
tags: [backend, api, auth, feature, refactor]
created: 2026-04-03
---

# API Backend Redesign — Auth + Chat + History

## Overview

Full clean-slate rebuild of the existing Go (Echo + pgx/PostgreSQL) backend.
All smartwatch IoT / TCP / Google OAuth code is **removed**.
Three features delivered:
1. **Auth** — email + password (argon2id), JWT access (15 min) + DB-backed refresh (30 d)
2. **Chat** — receive message → call OpenResponses `POST /v1/responses` → save Q&A → return answer
3. **History** — paginated message list for the authenticated user

## Brainstorm Reference

[`plans/reports/brainstorm-260403-0933-api-backend-redesign.md`](../reports/brainstorm-260403-0933-api-backend-redesign.md)

## Phases

| # | Phase | Status | Effort | Link |
|---|-------|--------|--------|------|
| 1 | Cleanup — remove old code & deps | Pending | 30m | [phase-01](./phase-01-cleanup.md) |
| 2 | Foundations — config, DB, migrations, domain | Pending | 45m | [phase-02-foundations.md](./phase-02-foundations.md) |
| 3 | Auth — register/login/refresh/logout | Pending | 2h | [phase-03-auth.md](./phase-03-auth.md) |
| 4 | Chat — message endpoint + OpenResponses | Pending | 1.5h | [phase-04-chat.md](./phase-04-chat.md) |
| 5 | History — paginated list | Pending | 45m | [phase-05-history.md](./phase-05-history.md) |
| 6 | Wiring — main.go, router, Makefile | Pending | 30m | [phase-06-wiring.md](./phase-06-wiring.md) |

## Key Dependencies

- `golang.org/x/crypto` — argon2id (already in go.mod)
- `github.com/golang-jwt/jwt/v5` — **add** via `go get`
- `golang.org/x/oauth2` — **remove** from go.mod
- PostgreSQL running locally or via `DATABASE_URL`
