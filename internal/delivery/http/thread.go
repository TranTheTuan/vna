package http

import (
	"github.com/labstack/echo/v4"

	"github.com/TranTheTuan/vna/configs"
	http_handler "github.com/TranTheTuan/vna/internal/handler/http"
)

// RegisterThreadRoutes registers thread routes under the given API group.
func RegisterThreadRoutes(e *echo.Group, th *http_handler.ThreadHandler, cfg *configs.Config) {
	threads := e.Group("/threads")
	threads.Use(JWTMiddleware(cfg))
	threads.GET("", th.List)
	threads.PUT("/:id", th.Rename)
}
