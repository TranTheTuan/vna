package http

import (
	"github.com/labstack/echo/v4"

	"github.com/TranTheTuan/vna/configs"
	http_handler "github.com/TranTheTuan/vna/internal/handler/http"
)

func RegisterAuthRoutes(e *echo.Group, ah *http_handler.AuthHandler, cfg *configs.Config) {
	// Public auth routes
	auth := e.Group("/auth")
	auth.POST("/register", ah.Register)
	auth.POST("/login", ah.Login)
	auth.POST("/refresh", ah.Refresh)

	// Protected auth route (logout requires a valid access token)
	authProtected := e.Group("/auth")
	authProtected.Use(JWTMiddleware(cfg))
	authProtected.POST("/logout", ah.Logout)
}
