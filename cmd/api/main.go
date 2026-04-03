// @title           VNA API
// @version         1.0
// @description     VNA backend API — authentication and AI chat.
// @host            localhost:8080
// @BasePath        /
// @securityDefinitions.apikey BearerAuth
// @in              header
// @name            Authorization
// @description     Enter: Bearer {token}

package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"

	"github.com/TranTheTuan/vna/configs"
	"github.com/TranTheTuan/vna/internal/db"
	http_delivery "github.com/TranTheTuan/vna/internal/delivery/http"
	http_handler "github.com/TranTheTuan/vna/internal/handler/http"
	"github.com/TranTheTuan/vna/internal/repository"
	"github.com/TranTheTuan/vna/internal/service"
)

func main() {
	cfg, err := configs.LoadConfig()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	// Structured JSON logger — writes to stdout, level Info
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))

	pool, err := db.NewPool(cfg.Database.BuildConnectionString())
	if err != nil {
		log.Fatalf("db: %v", err)
	}
	defer pool.Close()

	// Repositories
	userRepo := repository.NewUserRepository(pool)
	messageRepo := repository.NewMessageRepository(pool)

	// Services
	userSvc := service.NewUserService(cfg, userRepo, logger)
	messageSvc := service.NewMessageService(cfg, messageRepo, logger)

	// Handlers
	authHandler := http_handler.NewAuthHandler(userSvc)
	messageHandler := http_handler.NewMessageHandler(messageSvc)

	// Echo server
	e := echo.New()
	e.HideBanner = true
	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())

	// Swagger UI — only enabled when SWAGGER_ENABLED=true
	if cfg.SwaggerEnabled {
		e.GET("/swagger/*", echoSwagger.WrapHandler)
	}

	apiGroup := e.Group("/api/v1")
	http_delivery.RegisterAuthRoutes(apiGroup, authHandler, cfg)
	http_delivery.RegisterMessageRoutes(apiGroup, messageHandler, cfg)

	// Graceful shutdown on SIGINT / SIGTERM
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		<-ctx.Done()
		log.Println("Shutting down...")
		if err := e.Shutdown(context.Background()); err != nil {
			log.Printf("shutdown error: %v", err)
		}
	}()

	log.Println("HTTP server listening on :8080")
	if err := e.Start(":8080"); err != nil {
		log.Println("Server stopped:", err)
	}
}
