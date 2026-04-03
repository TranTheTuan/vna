# Phase 01 — Cleanup

**Priority:** P1 | **Status:** Pending | **Effort:** 30m

## Overview

Remove all IoT/TCP/Google OAuth code and unused dependencies. Leave only the skeleton needed for the new features.

## Context Links

- Brainstorm: `plans/reports/brainstorm-260403-0933-api-backend-redesign.md`

## Files to DELETE

```
internal/domain/device.go
internal/domain/packet.go
internal/dto/device_dto.go
internal/dto/packet_dto.go
internal/demo/                          (entire dir)
cmd/hdp_app/                            (entire dir)
internal/delivery/tcp/                  (entire dir, if exists)
internal/handler/tcp/                   (entire dir, if exists)
internal/auth/                          (entire dir — Google OAuth helpers)
internal/api/                           (entire dir — template renderer)
internal/app/                           (empty dir)
internal/migrations/001_initial_schema.up.sql
internal/migrations/002_add_imei_suffix_index.up.sql
pkg/_your_public_lib_/                  (placeholder dir)
web/                                    (templates, if exists — no longer needed)
```

## Files to KEEP (do not touch yet)

```
cmd/api/main.go           ← rewritten in Phase 06
configs/config.go         ← rewritten in Phase 02
internal/db/connection.go ← unchanged
internal/delivery/http/router.go        ← rewritten in Phase 06
internal/handler/http/auth.go           ← rewritten in Phase 03
internal/repository/user.go            ← rewritten in Phase 03
internal/repository/message.go         ← rewritten in Phase 04
internal/service/user.go               ← rewritten in Phase 03
internal/service/message.go            ← rewritten in Phase 04
internal/dto/auth_dto.go               ← rewritten in Phase 03
internal/domain/user.go                ← extended in Phase 02
go.mod / go.sum
Makefile
```

## Implementation Steps

1. Delete files/dirs listed above using `rm -rf` (verify each path exists first)
2. Remove `golang.org/x/oauth2` from go.mod:
   ```bash
   go get golang.org/x/oauth2@none
   go mod tidy
   ```
3. Add new JWT dependency:
   ```bash
   go get github.com/golang-jwt/jwt/v5
   go mod tidy
   ```
4. Verify project still compiles (will have broken imports — that's expected, fixed in later phases):
   ```bash
   go build ./... 2>&1 | head -30
   ```

## Todo List

- [ ] Delete all IoT/TCP/OAuth files and dirs
- [ ] Remove `golang.org/x/oauth2` from go.mod
- [ ] Add `github.com/golang-jwt/jwt/v5` to go.mod
- [ ] Run `go mod tidy`

## Success Criteria

- `go.mod` no longer references `golang.org/x/oauth2` or `cloud.google.com/go/compute/metadata`
- `go.mod` references `github.com/golang-jwt/jwt/v5`
- All listed dirs/files are gone
