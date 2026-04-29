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

// Healthz godoc
// @Summary Liveness probe
// @Description Returns 200 if the process is alive
// @Tags health
// @Produce json
// @Success 200 {object} map[string]string
// @Router /healthz [get]
func healthzHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
}

// Readyz godoc
// @Summary Readiness probe
// @Description Returns 200 if DB is reachable, 503 otherwise
// @Tags health
// @Produce json
// @Success 200 {object} map[string]string
// @Failure 503 {object} map[string]string
// @Router /readyz [get]
func readyzHandler(pool *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		if err := pool.Ping(); err != nil {
			return c.JSON(http.StatusServiceUnavailable, map[string]string{"status": "error"})
		}
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	}
}
