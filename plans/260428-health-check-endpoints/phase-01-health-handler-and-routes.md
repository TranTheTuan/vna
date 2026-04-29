# Phase 01 — Health Handler & Routes

**Context:** `plans/260428-health-check-endpoints/plan.md`
**Priority:** P2 | **Status:** Complete | **Effort:** 0.5h

## Overview

Create health check route registration with two handlers in `internal/delivery/http/health.go`.

## Requirements

- `GET /healthz` — liveness: always returns 200 with `{"status":"ok"}`, no dependencies
- `GET /readyz` — readiness: pings DB pool via `Ping()`, returns 200 if ok, 503 if error
- Both endpoints public (no JWT middleware)
- Swagger annotations on both handlers

## Related Code Files

| Action | File |
|--------|------|
| Create | `internal/delivery/http/health.go` |

## Implementation Steps

1. Create `internal/delivery/http/health.go`:

```go
package http

import (
	"database/sql"
	"net/http"

	"github.com/labstack/echo/v4"
)

// RegisterHealthRoutes registers /healthz and /readyz at root level (no JWT).
func RegisterHealthRoutes(e *echo.Echo, pool *sql.DB) {
	e.GET("/healthz", healthzHandler)
	e.GET("/readyz", readyzHandler(pool))
}

// healthzHandler always returns 200 — process is alive.
func healthzHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
}

// readyzHandler returns 200 if DB ping succeeds, 503 otherwise.
func readyzHandler(pool *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		if err := pool.Ping(); err != nil {
			return c.JSON(http.StatusServiceUnavailable, map[string]string{"status": "error"})
		}
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	}
}
```

2. Add Swagger annotations:

```go
// Healthz godoc
// @Summary Liveness probe
// @Description Returns 200 if the process is alive
// @Tags health
// @Produce json
// @Success 200 {object} map[string]string
// @Router /healthz [get]
```

```go
// Readyz godoc
// @Summary Readiness probe
// @Description Returns 200 if DB is reachable, 503 otherwise
// @Tags health
// @Produce json
// @Success 200 {object} map[string]string
// @Failure 503 {object} map[string]string
// @Router /readyz [get]
```

## Todo

- [ ] Create `internal/delivery/http/health.go`
- [ ] Add `RegisterHealthRoutes` function
- [ ] Implement `healthzHandler` (always 200)
- [ ] Implement `readyzHandler` (pings DB pool)
- [ ] Add Swagger annotations to both handlers

## Success Criteria

- `RegisterHealthRoutes` accepts `*echo.Echo` and `*sql.DB`
- `healthzHandler` returns 200 with `{"status":"ok"}`
- `readyzHandler` returns 200/503 based on DB ping result
- Endpoints at root level (`/healthz`, `/readyz`), not under `/api/v1`

## Risk Assessment

- Minimal risk — standard pattern, follows KISS
- Using `*sql.DB.Ping()` — lightweight, no additional deps
