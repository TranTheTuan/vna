---
title: "Phase 2: Register Middleware and Expose /metrics"
phase: 2
status: completed
effort: 0.5h
---

# Phase 2: Register Middleware and Expose /metrics

## Objective
Register the Prometheus metrics middleware and expose the `/metrics` endpoint for scraping.

## Tasks
1. **Update imports in `cmd/api/main.go`**
   ```go
   import (
       // ... existing imports ...
       "github.com/globocom/echo-prometheus"
       "github.com/prometheus/client_golang/prometheus/promhttp"
   )
   ```

2. **Register metrics middleware** (after existing middleware, before routes)
   ```go
   e.Use(echoPrometheus.MetricsMiddleware())
   ```
   - Place after `middleware.CORSWithConfig(...)` and before route registration
   - Captures request duration, count, and status codes by route

3. **Expose `/metrics` endpoint** (public, no JWT, similar to health endpoints)
   ```go
   e.GET("/metrics", echo.WrapHandler(promhttp.Handler()))
   ```
   - Place after health routes (`http_delivery.RegisterHealthRoutes(e, pool)`)
   - No auth middleware — standard for Prometheus scraping
   - Uses `echo.WrapHandler` to adapt `promhttp.Handler()` to Echo

4. **Update Swagger annotation** in `cmd/api/main.go`
   ```go
   // @title           VNA API
   // @description     VNA backend API — authentication, AI chat, and Prometheus metrics.
   ```
   - Add mention of Prometheus metrics to description

## Files Modified
- `cmd/api/main.go`

## Verification
- [ ] `go build ./cmd/api/...` succeeds
- [ ] `curl http://localhost:8080/metrics` returns Prometheus-formatted metrics
- [ ] Metrics include `echo_request_duration_seconds` and `echo_request_total`
- [ ] Making API requests increases metric counters
