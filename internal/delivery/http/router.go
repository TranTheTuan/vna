// Package http registers all HTTP routes and middleware for the Echo server.
package http

import (
	"strings"

	"github.com/labstack/echo/v4"

	"github.com/TranTheTuan/vna/configs"
	http_handler "github.com/TranTheTuan/vna/internal/handler/http"
	"github.com/TranTheTuan/vna/pkg/jwt_util"

	_ "github.com/TranTheTuan/vna/internal/docs"
)

// RegisterRoutes wires all handlers and middleware to the Echo instance.
func RegisterRoutes(e *echo.Echo, cfg *configs.Config, ah *http_handler.AuthHandler, mh *http_handler.MessageHandler) {
	// Public auth routes
	auth := e.Group("/api/v1/auth")
	auth.POST("/register", ah.Register)
	auth.POST("/login", ah.Login)
	auth.POST("/refresh", ah.Refresh)

	// Protected auth route (logout requires a valid access token)
	authProtected := e.Group("/api/v1/auth")
	authProtected.Use(JWTMiddleware(cfg))
	authProtected.POST("/logout", ah.Logout)

	// Protected message routes
	msgs := e.Group("/api/v1/messages")
	msgs.Use(JWTMiddleware(cfg))
	msgs.POST("", mh.Send)
	msgs.GET("", mh.List)
}

// JWTMiddleware validates Bearer access tokens and injects user claims into context.
func JWTMiddleware(cfg *configs.Config) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			header := c.Request().Header.Get("Authorization")
			if !strings.HasPrefix(header, "Bearer ") {
				return echo.ErrUnauthorized
			}
			tokenStr := strings.TrimPrefix(header, "Bearer ")

			claims, err := jwt_util.ParseAccessToken(tokenStr, cfg.Auth.JWTSecret)
			if err != nil {
				return echo.ErrUnauthorized
			}

			c.Set("user_id", claims.UserID)
			c.Set("user_email", claims.Email)
			return next(c)
		}
	}
}
