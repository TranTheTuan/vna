# Brainstorm Report: API Backend Redesign
**Date:** 2026-04-03 | **Status:** Agreed

---

## Problem Statement

Existing project is a Go (Echo + pgx/PostgreSQL) IoT/smartwatch backend with Google OAuth, TCP server, device management. Requirements have shifted to a **chat API backend** with:
1. Email/password register & login → return JWT pair
2. Receive message → call `POST /v1/responses` (OpenResponses, OpenAI-compatible) → save Q&A → return answer
3. List message history for logged-in user (paginated)

All old IoT code (TCP server, devices, packets, Google OAuth) is **removed**. Same Go repo, clean slate.

---

## Key Decisions (Agreed)

| Decision | Choice | Rationale |
|---|---|---|
| Auth method | Email + password | New requirement; replaces Google OAuth |
| Password hashing | **argon2id** | Gold standard; `golang.org/x/crypto` already in go.mod |
| Token strategy | Access + Refresh pair | DB-backed refresh → revocation support |
| External API | OpenResponses `POST /v1/responses` | OpenAI Responses API compatible |
| Message history | Paginated (cursor/offset) | Scalability |
| Refresh token storage | DB (hash stored) | Enables logout/revocation |
| Makefile | Updated | Convenience |

---

## Architecture

### Layer Structure (Clean Architecture, existing pattern)
```
cmd/api/main.go           ← entrypoint, DI wiring
configs/config.go         ← env config (updated)
internal/
  domain/                 ← pure data models (user.go, message.go)
  dto/                    ← request/response shapes (auth_dto.go, message_dto.go)
  db/connection.go        ← unchanged
  migrations/             ← new SQL migrations
  repository/             ← user.go (w/ refresh token), message.go
  service/                ← user.go (auth logic), message.go (AI call + save)
  handler/http/           ← auth.go, message.go
  delivery/http/router.go ← updated routes + JWT middleware
pkg/                      ← jwt util, argon2 util (new small helpers)
```

### Files to REMOVE
- `internal/domain/device.go`, `packet.go`
- `internal/dto/device_dto.go`, `packet_dto.go`
- `internal/demo/` (all)
- All `internal/handler/http/` except auth.go (rewrite it)
- `cmd/hdp_app/` (entire directory)
- Any TCP-related code in delivery, handler, service
- `internal/app/` (empty, remove)
- `internal/migrations/001_initial_schema.up.sql`, `002_add_imei_suffix_index.up.sql`
- All OAuth-related internal packages (`internal/auth/` if exists, `internal/api/`)

### Files to KEEP & UPDATE
- `cmd/api/main.go` → complete rewrite (remove TCP, OAuth, devices)
- `configs/config.go` → swap Google OAuth fields for JWT secret, OpenResponses API URL/key
- `internal/db/connection.go` → unchanged (pgx pool)
- `internal/domain/user.go` → extend with password hash + created_at
- `internal/handler/http/auth.go` → complete rewrite for email/password + JWT
- `internal/delivery/http/router.go` → update routes
- `Makefile` → update targets

### Files to CREATE
- `internal/domain/message.go` → Message struct (id, user_id, question, answer, created_at)
- `internal/dto/message_dto.go` → request/response DTOs
- `internal/repository/user.go` → CreateUser, FindByEmail, SaveRefreshToken, RevokeRefreshToken
- `internal/repository/message.go` → SaveMessage, ListByUser(cursor, limit)
- `internal/service/user.go` → Register, Login, RefreshToken, Logout
- `internal/service/message.go` → SendMessage (calls OpenResponses + saves)
- `internal/handler/http/message.go` → SendMessage handler, ListHistory handler
- `internal/migrations/001_users.up.sql` → users + refresh_tokens tables
- `internal/migrations/002_messages.up.sql` → messages table with index
- `pkg/jwtutil/jwt.go` → GenerateAccessToken, GenerateRefreshToken, ParseToken
- `pkg/argon2util/argon2.go` → HashPassword, VerifyPassword

---

## API Design

### Auth
```
POST /api/v1/auth/register    { email, password } → { user_id, email }
POST /api/v1/auth/login       { email, password } → { access_token, refresh_token, expires_in }
POST /api/v1/auth/refresh     { refresh_token }   → { access_token, expires_in }
POST /api/v1/auth/logout      Bearer token        → 200 OK (revokes refresh token)
```

### Messages (all require JWT Bearer auth)
```
POST /api/v1/messages         { message: "..." }  → { id, question, answer, created_at }
GET  /api/v1/messages         ?limit=20&cursor=<id> → { data: [...], next_cursor }
```

---

## Database Schema

### `users` table
```sql
id            UUID PRIMARY KEY DEFAULT gen_random_uuid()
email         TEXT NOT NULL UNIQUE
password_hash TEXT NOT NULL           -- argon2id encoded
created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
```

### `refresh_tokens` table
```sql
id          UUID PRIMARY KEY DEFAULT gen_random_uuid()
user_id     UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE
token_hash  TEXT NOT NULL UNIQUE    -- SHA-256 of the raw token
expires_at  TIMESTAMPTZ NOT NULL
created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
revoked_at  TIMESTAMPTZ             -- NULL = valid
```

### `messages` table
```sql
id          UUID PRIMARY KEY DEFAULT gen_random_uuid()
user_id     UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE
question    TEXT NOT NULL
answer      TEXT NOT NULL
created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
```
**Index:** `(user_id, created_at DESC)` for efficient history pagination.

---

## JWT Strategy

- **Access token**: HS256, 15-min TTL, claims: `{sub: user_id, email, exp}`
- **Refresh token**: opaque random 32-byte token, stored as SHA-256 hash in DB, 30-day TTL
- Middleware extracts Bearer token from `Authorization` header, parses/validates JWT
- Refresh endpoint: verify token hash in DB (not revoked, not expired) → issue new access token

---

## OpenResponses API Integration

Request to `POST /v1/responses` (OpenAI Responses API compatible):
```json
{
  "model": "<from config>",
  "input": [{ "role": "user", "content": "<user message>" }]
}
```
Headers: `Authorization: Bearer <API_KEY>`, `Content-Type: application/json`

Response: extract `output[0].content[0].text` as the answer string.

All handled in `service/message.go` — no streaming for now (KISS).

---

## Config Changes

```
OLD (remove): GOOGLE_CLIENT_ID, GOOGLE_CLIENT_SECRET, GOOGLE_REDIRECT_URL, SESSION_SECRET, TCP_ADDR
NEW (add):    JWT_SECRET, JWT_ACCESS_TTL (15m), JWT_REFRESH_TTL (720h)
              OPENRESPONSES_API_URL, OPENRESPONSES_API_KEY, OPENRESPONSES_MODEL
```

---

## New Dependencies Needed
- `github.com/golang-jwt/jwt/v5` — JWT generation/parsing
- No new dep for argon2 (`golang.org/x/crypto` already there)
- No new dep for password hashing

---

## What to Remove from go.mod
- `golang.org/x/oauth2` — no longer needed
- `cloud.google.com/go/compute/metadata` — indirect, pulled by oauth2

---

## Risk Assessment

| Risk | Mitigation |
|---|---|
| OpenResponses API timeout | Set HTTP timeout (10s), return 504 to client |
| Argon2 slow on login | Acceptable; it's designed to be slow. Only on auth path. |
| Refresh token theft | HTTPS only, httpOnly cookie optional, 30-day TTL, revocable |
| DB migration ordering | Use numbered up-migrations only; no down needed for now |

---

## Success Criteria
- `POST /api/v1/auth/register` creates user with argon2id hashed password
- `POST /api/v1/auth/login` returns valid JWT pair
- `POST /api/v1/messages` calls OpenResponses, saves both Q&A, returns answer
- `GET /api/v1/messages` returns paginated history for authed user only
- All IoT/TCP code removed, project compiles cleanly
- Makefile: `make run`, `make build`, `make migrate`

---

## Unresolved Questions
- Does the OpenResponses API at openclaw.ai return streaming SSE or full JSON response? (Assumed full JSON for now — KISS)
- Should register validate email format server-side? (Recommend: yes, simple regex or stdlib check)
- Should old migrations be deleted or just replaced by new ones? (Recommend: delete old, create fresh numbered from 001)
