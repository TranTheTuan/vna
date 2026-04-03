# Phase 06 — Wiring (main.go, Router, Makefile)

**Priority:** P1 | **Status:** Pending | **Effort:** 30m

## Overview

Wire everything together: DI in main.go, route registration with JWT middleware, Makefile targets.

## Context Links

- Depends on all previous phases

## Files to REWRITE

### `cmd/api/main.go`

```go
package main

import (
    "context"
    "log"
    "os"
    "os/signal"
    "syscall"

    "github.com/labstack/echo/v4"
    "github.com/labstack/echo/v4/middleware"

    "github.com/TranTheTuan/vna/configs"
    "github.com/TranTheTuan/vna/internal/db"
    http_delivery "github.com/TranTheTuan/vna/internal/delivery/http"
    http_handler "github.com/TranTheTuan/vna/internal/handler/http"
    "github.com/TranTheTuan/vna/internal/repository"
    "github.com/TranTheTuan/vna/internal/service"
)

func main() {
    cfg := configs.LoadConfig()

    pool, err := db.NewPool(cfg.DatabaseURL)
    if err != nil { log.Fatalf("db: %v", err) }
    defer pool.Close()

    // Repositories
    userRepo    := repository.NewUserRepository(pool)
    messageRepo := repository.NewMessageRepository(pool)

    // Services
    userSvc    := service.NewUserService(cfg, userRepo)
    messageSvc := service.NewMessageService(cfg, messageRepo)

    // Handlers
    authHandler    := http_handler.NewAuthHandler(cfg, userSvc)
    messageHandler := http_handler.NewMessageHandler(messageSvc)

    // Echo
    e := echo.New()
    e.Use(middleware.Logger())
    e.Use(middleware.Recover())

    http_delivery.RegisterRoutes(e, cfg, authHandler, messageHandler)

    // Graceful shutdown
    ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
    defer stop()

    go func() {
        <-ctx.Done()
        e.Shutdown(context.Background())
    }()

    log.Println("HTTP server starting on :8080")
    if err := e.Start(":8080"); err != nil {
        log.Println("Server stopped:", err)
    }
}
```

### `internal/delivery/http/router.go`

```go
package http

import (
    "github.com/labstack/echo/v4"
    "github.com/TranTheTuan/vna/configs"
    http_handler "github.com/TranTheTuan/vna/internal/handler/http"
)

func RegisterRoutes(
    e *echo.Echo,
    cfg *configs.Config,
    ah *http_handler.AuthHandler,
    mh *http_handler.MessageHandler,
) {
    // Auth (public)
    auth := e.Group("/api/v1/auth")
    auth.POST("/register", ah.Register)
    auth.POST("/login",    ah.Login)
    auth.POST("/refresh",  ah.Refresh)
    auth.POST("/logout",   ah.Logout)

    // Messages (protected)
    msgs := e.Group("/api/v1/messages")
    msgs.Use(JWTMiddleware(cfg))
    msgs.POST("",  mh.Send)
    msgs.GET("",   mh.List)
}
```

`JWTMiddleware` function (in same file):
```go
func JWTMiddleware(cfg *configs.Config) echo.MiddlewareFunc {
    return func(next echo.HandlerFunc) echo.HandlerFunc {
        return func(c echo.Context) error {
            header := c.Request().Header.Get("Authorization")
            if !strings.HasPrefix(header, "Bearer ") {
                return echo.ErrUnauthorized
            }
            tokenStr := strings.TrimPrefix(header, "Bearer ")
            claims, err := jwtutil.ParseAccessToken(tokenStr, cfg.JWTSecret)
            if err != nil {
                return echo.ErrUnauthorized
            }
            c.Set("user_id",    claims.UserID)
            c.Set("user_email", claims.Email)
            return next(c)
        }
    }
}
```

### `Makefile`

```makefile
.PHONY: run build migrate tidy

run:
	go run ./cmd/api/...

build:
	go build -o bin/api ./cmd/api/...

migrate:
	@for f in internal/migrations/*.up.sql; do \
		echo "Applying $$f..."; \
		psql "$$DATABASE_URL" -f "$$f"; \
	done

tidy:
	go mod tidy
```

## Environment Variables Reference (`.env` template)

Document in a `.env.example` (create new file):
```
DATABASE_URL=postgres://user:pass@localhost:5432/dbname?sslmode=disable
JWT_SECRET=change-me-to-a-long-random-string
JWT_ACCESS_TTL=15m
JWT_REFRESH_TTL=720h
OPENRESPONSES_URL=https://api.openclaw.ai
OPENRESPONSES_API_KEY=your-api-key-here
OPENRESPONSES_MODEL=gpt-4o
```

## Todo List

- [ ] Rewrite `cmd/api/main.go` (DI wiring, no TCP)
- [ ] Rewrite `internal/delivery/http/router.go` (routes + JWT middleware)
- [ ] Update `Makefile` with `run`, `build`, `migrate`, `tidy` targets
- [ ] Create `.env.example`
- [ ] `go build ./...` compiles with zero errors
- [ ] `go vet ./...` passes
- [ ] Run `make migrate` against local DB
- [ ] Smoke test: register → login → send message → list history

## Success Criteria

- `make build` produces `bin/api` binary
- `make run` starts server on :8080
- `make migrate` applies both SQL files without errors
- Full flow works end-to-end:
  1. `POST /api/v1/auth/register` → 201
  2. `POST /api/v1/auth/login` → 200 with tokens
  3. `POST /api/v1/messages` (with Bearer token) → 201 with answer
  4. `GET /api/v1/messages` (with Bearer token) → 200 with paginated list
