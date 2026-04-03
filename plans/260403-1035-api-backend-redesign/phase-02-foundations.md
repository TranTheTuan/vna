# Phase 02 — Foundations

**Priority:** P1 | **Status:** Pending | **Effort:** 45m

## Overview

Update config, extend domain models, create DB migrations, add `pkg` utilities for argon2id and JWT.

## Context Links

- Brainstorm: `plans/reports/brainstorm-260403-0933-api-backend-redesign.md`
- Phase 01 must complete first (cleanup + deps)

## Files to MODIFY

### `configs/config.go`

Replace Google OAuth / TCP fields with:

```go
type Config struct {
    DatabaseURL          string
    JWTSecret            string
    JWTAccessTTL         time.Duration  // default: 15m
    JWTRefreshTTL        time.Duration  // default: 720h (30d)
    OpenResponsesURL     string         // e.g. https://api.openclaw.ai
    OpenResponsesAPIKey  string
    OpenResponsesModel   string         // e.g. gpt-4o
}
```

Env vars to read:
```
DATABASE_URL
JWT_SECRET
JWT_ACCESS_TTL      (default: "15m")
JWT_REFRESH_TTL     (default: "720h")
OPENRESPONSES_URL
OPENRESPONSES_API_KEY
OPENRESPONSES_MODEL (default: "gpt-4o")
```

Add helper `parseDuration(key, fallback string) time.Duration`.

### `internal/domain/user.go`

```go
type User struct {
    ID           string    // UUID
    Email        string
    PasswordHash string    // argon2id encoded string
    CreatedAt    time.Time
}
```

## Files to CREATE

### `internal/domain/message.go`

```go
type Message struct {
    ID        string
    UserID    string
    Question  string
    Answer    string
    CreatedAt time.Time
}
```

### `pkg/argon2_util/argon2_util.go`

Two exported functions:
- `HashPassword(password string) (string, error)` — argon2id, encode as `$argon2id$v=19$m=65536,t=3,p=4$<salt>$<hash>` (PHC string format)
- `VerifyPassword(password, encoded string) (bool, error)` — parse PHC, re-hash, compare

Use `golang.org/x/crypto/argon2`. Parameters: memory=65536, time=3, threads=4, keyLen=32, saltLen=16.

### `pkg/jwt_util/jwt_util.go`

Three exported functions:
- `GenerateAccessToken(userID, email, secret string, ttl time.Duration) (string, error)` — HS256, claims: `sub`, `email`, `exp`, `iat`
- `GenerateRefreshToken() (raw string, hash string, err error)` — 32 random bytes, base64url raw, SHA-256 hex hash
- `ParseAccessToken(tokenStr, secret string) (*Claims, error)` — returns custom `Claims` struct with `UserID`, `Email`

```go
type Claims struct {
    UserID string
    Email  string
    jwt.RegisteredClaims
}
```

## DB Migrations

### `internal/migrations/001_users_and_refresh_tokens.up.sql`

```sql
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE IF NOT EXISTS users (
    id            UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    email         TEXT        NOT NULL UNIQUE,
    password_hash TEXT        NOT NULL,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS refresh_tokens (
    id         UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id    UUID        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_hash TEXT        NOT NULL UNIQUE,  -- SHA-256 hex of raw token
    expires_at TIMESTAMPTZ NOT NULL,
    revoked_at TIMESTAMPTZ,                  -- NULL = valid
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_refresh_tokens_user
    ON refresh_tokens(user_id);
```

### `internal/migrations/002_messages.up.sql`

```sql
CREATE TABLE IF NOT EXISTS messages (
    id         UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id    UUID        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    question   TEXT        NOT NULL,
    answer     TEXT        NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Efficient paginated history per user
CREATE INDEX IF NOT EXISTS idx_messages_user_time
    ON messages(user_id, created_at DESC);
```

## Todo List

- [ ] Rewrite `configs/config.go` with new fields and env vars
- [ ] Extend `internal/domain/user.go` with `PasswordHash`, `CreatedAt`
- [ ] Create `internal/domain/message.go`
- [ ] Create `pkg/argon2_util/argon2_util.go` with Hash/Verify
- [ ] Create `pkg/jwt_util/jwt_util.go` with Generate/Parse functions
- [ ] Create `internal/migrations/001_users_and_refresh_tokens.up.sql`
- [ ] Create `internal/migrations/002_messages.up.sql`
- [ ] Run `go build ./...` — should compile (no handler wiring yet)

## Success Criteria

- `go build ./pkg/...` compiles cleanly
- Config loads all new env vars with sensible defaults
- Migration SQL is valid (test with `psql -f migration.sql`)

## Security Considerations

- argon2id params (m=65536, t=3, p=4) are OWASP-recommended minimums
- JWT secret loaded from env — never hardcoded
- Refresh token stored as SHA-256 hash only — raw token never persisted
