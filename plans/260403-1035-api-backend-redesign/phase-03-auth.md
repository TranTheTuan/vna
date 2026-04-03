# Phase 03 — Auth (Register / Login / Refresh / Logout)

**Priority:** P1 | **Status:** Pending | **Effort:** 2h

## Overview

Implement full email+password auth: repository layer → service layer → HTTP handler → routes.

## Context Links

- Brainstorm: `plans/reports/brainstorm-260403-0933-api-backend-redesign.md`
- Depends on Phase 02 (domain, migrations, pkg utils must exist)

## API Contract

```
POST /api/v1/auth/register
  Body:    { "email": "user@example.com", "password": "..." }
  Success: 201 { "user_id": "uuid", "email": "..." }
  Errors:  400 invalid email/password, 409 email already registered

POST /api/v1/auth/login
  Body:    { "email": "...", "password": "..." }
  Success: 200 { "access_token": "...", "refresh_token": "...", "expires_in": 900 }
  Errors:  400 bad request, 401 invalid credentials

POST /api/v1/auth/refresh
  Body:    { "refresh_token": "..." }
  Success: 200 { "access_token": "...", "expires_in": 900 }
  Errors:  401 invalid/expired/revoked token

POST /api/v1/auth/logout
  Header:  Authorization: Bearer <access_token>
  Body:    { "refresh_token": "..." }
  Success: 200 {}
  Errors:  401 unauthorized
```

## Files to REWRITE

### `internal/dto/auth_dto.go`

```go
// Request DTOs
type RegisterRequest  { Email, Password string }
type LoginRequest     { Email, Password string }
type RefreshRequest   { RefreshToken string `json:"refresh_token"` }
type LogoutRequest    { RefreshToken string `json:"refresh_token"` }

// Response DTOs
type AuthResponse     { AccessToken, RefreshToken string; ExpiresIn int }
type RegisterResponse { UserID, Email string }
```

### `internal/repository/user.go`

Interface + pgx implementation:

```go
type UserRepository interface {
    Create(ctx, email, passwordHash string) (*domain.User, error)
    FindByEmail(ctx, email string) (*domain.User, error)
    FindByID(ctx, id string) (*domain.User, error)
    SaveRefreshToken(ctx, userID, tokenHash string, expiresAt time.Time) error
    FindRefreshToken(ctx, tokenHash string) (*RefreshTokenRow, error)  // returns user_id, expires_at, revoked_at
    RevokeRefreshToken(ctx, tokenHash string) error
}
```

SQL queries:
- `INSERT INTO users(email, password_hash) VALUES($1,$2) RETURNING id, email, created_at`
- `SELECT id, email, password_hash, created_at FROM users WHERE email=$1`
- `SELECT id, email, created_at FROM users WHERE id=$1`
- `INSERT INTO refresh_tokens(user_id, token_hash, expires_at) VALUES($1,$2,$3)`
- `SELECT user_id, expires_at, revoked_at FROM refresh_tokens WHERE token_hash=$1`
- `UPDATE refresh_tokens SET revoked_at=NOW() WHERE token_hash=$1`

### `internal/service/user.go`

Interface + implementation:

```go
type UserService interface {
    Register(ctx, email, password string) (*domain.User, error)
    Login(ctx, email, password string) (accessToken, refreshToken string, err error)
    RefreshToken(ctx, rawRefreshToken string) (accessToken string, err error)
    Logout(ctx, rawRefreshToken string) error
}
```

Logic:
- **Register**: validate email (regexp `^[^@]+@[^@]+\.[^@]+$`), validate password len ≥ 8, hash with argon2id, call repo.Create; return 409 if duplicate email
- **Login**: repo.FindByEmail → argon2util.VerifyPassword → generate access+refresh tokens → repo.SaveRefreshToken (hash of raw)
- **RefreshToken**: SHA-256 hash the raw token → repo.FindRefreshToken → check not revoked + not expired → generate new access token
- **Logout**: SHA-256 hash → repo.RevokeRefreshToken

### `internal/handler/http/auth.go`

```go
type AuthHandler struct { svc service.UserService; cfg *configs.Config }
```

Four methods: `Register`, `Login`, `Refresh`, `Logout` — each:
1. Bind + validate request DTO
2. Call service method
3. Return JSON response or error

Error mapping:
- `service.ErrDuplicateEmail` → 409
- `service.ErrInvalidCredentials` → 401
- `service.ErrTokenExpired` / `ErrTokenRevoked` → 401
- All others → 500

### JWT Middleware (in `internal/delivery/http/router.go`)

```go
func JWTMiddleware(cfg *configs.Config) echo.MiddlewareFunc {
    return func(next echo.HandlerFunc) echo.HandlerFunc {
        return func(c echo.Context) error {
            // Extract Bearer token from Authorization header
            // jwtutil.ParseAccessToken → store userID + email in context
            // Return 401 on failure
        }
    }
}
```

Context keys: `"user_id"`, `"user_email"`.

## Todo List

- [ ] Rewrite `internal/dto/auth_dto.go`
- [ ] Implement `internal/repository/user.go` (interface + pgx impl)
- [ ] Implement `internal/service/user.go` (interface + impl with error sentinels)
- [ ] Rewrite `internal/handler/http/auth.go` (4 handlers)
- [ ] Add JWT middleware function (will be wired in Phase 06)

## Success Criteria

- `go build ./internal/...` compiles cleanly
- Register: duplicate email returns 409, weak password returns 400
- Login: wrong password returns 401, valid returns tokens
- Refresh: invalid/expired token returns 401
- Logout: revokes token, subsequent refresh returns 401

## Security Considerations

- Passwords never logged or returned in responses
- argon2id timing is constant (VerifyPassword always runs even if user not found, to prevent timing attacks)
- Refresh token stored as hash — raw token only in HTTP response
- Access token short-lived (15 min) to limit blast radius of theft
