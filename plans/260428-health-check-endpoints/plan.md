---
title: "Health Check Endpoints"
description: "Add /healthz and /readyz endpoints for liveness and readiness probes"
status: completed
priority: P2
effort: 1h
issue:
branch: feat/health-check
tags: [feature, backend, operations]
created: 2026-04-28
---

# Health Check Endpoints

## Overview

No health/readiness endpoints exist. Load balancers, K8s, and monitoring tools have no way to verify API status. Adds two standard endpoints:

- `GET /healthz` — liveness: always 200 `{"status":"ok"}`
- `GET /readyz` — readiness: 200 if DB ping succeeds, 503 otherwise

## Phases

| # | Phase | Status | Effort | Link |
|---|-------|--------|--------|------|
| 1 | Health handler & routes | Complete | 0.5h | [phase-01](./phase-01-health-handler-and-routes.md) |
| 2 | Wire-up & integration | Complete | 0.2h | [phase-02-wire-up-and-integration.md](./phase-02-wire-up-and-integration.md) |
| 3 | Swagger docs | Complete | 0.2h | [phase-03-swagger-docs.md](./phase-03-swagger-docs.md) |
| 4 | Tests & verification | Complete | 0.1h | [phase-04-tests-and-verification.md](./phase-04-tests-and-verification.md) |

## Dependencies

- Phase 1 → 2 → 3 → 4 (sequential)
- Phase 2 depends on DB pool being passed to health handler
- No external deps required
