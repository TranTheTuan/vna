package http

import (
	"github.com/labstack/echo/v4"

	"github.com/TranTheTuan/vna/configs"
	http_handler "github.com/TranTheTuan/vna/internal/handler/http"
)

func RegisterMessageRoutes(e *echo.Group, mh *http_handler.MessageHandler, cfg *configs.Config) {
	// Protected message routes
	msgs := e.Group("/messages")
	msgs.Use(JWTMiddleware(cfg))
	msgs.POST("", mh.Send)
	msgs.GET("", mh.List)
	msgs.POST("/stream", mh.SendStream)
}
