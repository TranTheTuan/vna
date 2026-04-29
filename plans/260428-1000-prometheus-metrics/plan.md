---
title: "Prometheus Metrics for VNA API"
description: "Add Prometheus metrics to VNA API using globocom/echo-prometheus middleware, expose /metrics endpoint for scraping, update Swagger documentation."
status: completed
priority: P2
effort: 1h
issue:
branch: feat/prometheus-metrics
tags: [feature, backend, observability]
created: 2026-04-28
---

# Prometheus Metrics for VNA API

## Overview
Add basic Prometheus HTTP metrics to the VNA API using the `globocom/echo-prometheus` middleware (compatible with Echo v4). Metrics include request count, latency, and status codes by route, plus standard Go process metrics. The `/metrics` endpoint is public (no JWT) to allow Prometheus scraping, consistent with existing health check endpoints.

## Success Criteria
- `go get github.com/globocom/echo-prometheus` succeeds
- `GET /metrics` returns valid Prometheus-formatted metrics
- Prometheus can scrape the endpoint (K8s service monitor ready)
- Swagger docs updated
- Build passes with no errors

## Phases
| # | Phase | Status | Effort | Link |
|---|-------|--------|--------|------|
| 1 | Add dependencies | completed | 0.25h | [phase-01-add-dependencies.md](./phase-01-add-dependencies.md) |
| 2 | Register middleware and expose /metrics | completed | 0.5h | [phase-02-middleware-and-endpoint.md](./phase-02-middleware-and-endpoint.md) |
| 3 | Documentation and verification | completed | 0.25h | [phase-03-docs-and-verify.md](./phase-03-docs-and-verify.md) |

## Dependencies
- Phases 1 → 2 → 3 (strict sequential)
- Phase 1 (dependencies) must complete before Phase 2 (code changes require imported packages)
- Phase 2 (code changes) must complete before Phase 3 (docs and verification depend on working code)
