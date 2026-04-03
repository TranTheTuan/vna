# Brainstorm Report: Swagger/OpenAPI Documentation for VNA Go API

**Date:** 2026-04-03  
**Scope:** Adding API documentation (Swagger UI) to `github.com/TranTheTuan/vna` (Echo v4)

---

## Problem Statement

The VNA project has 6 REST endpoints (auth + messages) with no machine-readable API contract or interactive explorer. Need to add Swagger/OpenAPI docs with minimal maintenance overhead for a solo project.

---

## Reference Project Scout (`sotabox_be`)

| Finding | Detail |
|---|---|
| Tool | `github.com/swaggo/swag v1.16.6` (code-annotation generator) |
| UI library | `github.com/iris-contrib/swagger` (Iris adapter — not reusable here) |
| Annotation location | Global annotations in `routers/swagger.go`; per-endpoint annotations directly above each handler function |
| UI route | `/swagger` + `/swagger/{any}`, protected with HTTP BasicAuth |
| Generated artifact | `src/docs/docs.go` (auto-generated, committed to repo) |
| `swag init` entrypoint | Scans all `.go` files under `src/` |

**Key takeaway:** Reference project uses the canonical `swaggo/swag` annotation approach. The only adapter difference is Iris vs Echo — for Echo the equivalent is `github.com/swaggo/echo-swagger`.

---

## Approaches Evaluated

### A) swaggo/swag + echo-swagger ✅ RECOMMENDED

- Write `// @Summary`, `// @Router`, `// @Param`, `// @Success` comments above each handler method
- Run `swag init -g cmd/api/main.go -o internal/docs` to generate `docs.go`
- Register one route: `e.GET("/swagger/*", echoSwagger.WrapHandler)`
- Artifacts: `internal/docs/docs.go`, `internal/docs/swagger.json`, `internal/docs/swagger.yaml`

**Pros:**
- Identical pattern to reference project — predictable, team-familiar
- Annotations live next to the code — harder to go stale
- Single `make docs` target regenerates everything
- `echo-swagger` is the de-facto Echo adapter (`github.com/swaggo/echo-swagger`)
- DTO structs (`dto.RegisterRequest`, etc.) are already clean — minimal extra annotation needed
- Only 6 endpoints + small DTO set → annotation effort is ~2–3 hours, one-time
- Generated `docs.go` can be committed so UI works without `swag` installed in CI

**Cons:**
- Requires `swag` CLI install (`go install github.com/swaggo/swag/cmd/swag@latest`)
- Generated `docs.go` is verbose boilerplate (but auto-managed)
- Annotation syntax is slightly quirky (space-sensitive, no IDE autocomplete out of the box)

---

### B) Hand-written OpenAPI YAML

- Write `api/openapi.yaml` by hand (OpenAPI 3.0)
- Serve via `e.GET("/swagger/*", echo.WrapHandler(redoc or rapidoc static handler))`

**Pros:**
- No codegen tooling dependency
- Full OpenAPI 3.0 (swaggo generates 2.0)
- Clean separation of spec from code

**Cons:**
- No enforcement that spec matches code — can drift silently
- More upfront effort for even a small API
- Must manually update on every field/route change
- YAGNI: overkill for 6 endpoints, solo project

---

## Recommendation: swaggo/swag + echo-swagger

Mirrors the reference project's proven pattern, ties docs to code via annotations, and is the lowest-friction path for a small Echo API.

### Packages to add

```
go get github.com/swaggo/swag           # annotation processor (CLI dep only)
go get github.com/swaggo/echo-swagger   # Echo UI handler
go get github.com/swaggo/files          # embedded swagger-ui assets
```

CLI tool (not a runtime dep):
```
go install github.com/swaggo/swag/cmd/swag@latest
```

### Where annotations go

| Location | What |
|---|---|
| `cmd/api/main.go` | Global: `@title`, `@version`, `@description`, `@host`, `@basePath`, `@securityDefinitions.apikey BearerAuth` |
| `internal/handler/http/auth.go` | Per-method: `@Summary`, `@Tags`, `@Accept`, `@Produce`, `@Param`, `@Success`, `@Failure`, `@Router` |
| `internal/handler/http/message.go` | Same pattern, including `@Security BearerAuth` for protected endpoints |

DTO structs in `internal/dto/` already have `json` tags — `swag` will pick them up automatically via `$ref` without extra annotation.

### UI route

Add to `internal/delivery/http/router.go`:
```go
import echoSwagger "github.com/swaggo/echo-swagger"

e.GET("/swagger/*", echoSwagger.WrapHandler)
```
Accessible at `http://localhost:8080/swagger/index.html`.

**Optional:** Gate behind env flag (`SWAGGER_ENABLED=true`) so it's disabled in production.

### Makefile target

```makefile
.PHONY: docs
docs:
    swag init -g cmd/api/main.go -o internal/docs --parseDependency --parseInternal
```

Add `docs` as prerequisite of `build` if you want always-fresh docs:
```makefile
build: docs
    go build -o bin/api ./cmd/api/...
```

### Echo v4 Gotchas

- `echoSwagger.WrapHandler` must be registered **before** any catch-all route or wildcard group
- Import the generated `docs` package as a blank import in `main.go` (`_ "github.com/TranTheTuan/vna/internal/docs"`) — this is the side-effect registration pattern swaggo uses
- `swag init` default entrypoint flag is `-g main.go`; with the `cmd/api/` layout, must specify `-g cmd/api/main.go`
- Echo's `/*` wildcard syntax: use `"/swagger/*"` not `"/swagger/{any}"` (that's Iris syntax)

---

## Implementation Steps

1. `go get` the two runtime packages (`echo-swagger`, `files`)
2. Install `swag` CLI
3. Add global annotations to `cmd/api/main.go` (8–10 lines)
4. Add per-handler annotations to `auth.go` (4 handlers × ~10 lines each)
5. Add per-handler annotations to `message.go` (2 handlers × ~10 lines each)
6. Run `make docs` → generates `internal/docs/`
7. Blank-import `internal/docs` in `main.go`
8. Register `/swagger/*` route in `router.go`
9. Verify UI at `localhost:8080/swagger/index.html`
10. Add `make docs` to `Makefile` and commit generated `docs/` (or `.gitignore` it and add CI step)

---

## Risks

| Risk | Mitigation |
|---|---|
| Annotations drift from code | `make docs` in CI; fail build if `docs.go` diff detected |
| `swag` CLI not on dev machines | Document in README; `go install` one-liner |
| Swagger UI exposed in prod | Env-gate the route registration |
| `swag` generates OpenAPI 2.0 (not 3.0) | Acceptable for this API size; upgrade path via `swag` v2 (still alpha) or future migration to hand-written YAML |

---

## Summary

For a 6-endpoint solo Go/Echo project, **swaggo/swag + echo-swagger** is the correct KISS choice. It mirrors the reference project, binds docs to code, and costs ~3h one-time effort. Hand-written YAML is overkill and drifts silently.

---

## Unresolved Questions

- Should the Swagger UI be env-gated (disabled in prod) or left open? → Recommend env flag `SWAGGER_ENABLED`
- Should `internal/docs/` be committed to git or regenerated in CI? → Either works; committing is simpler for solo project
- Any preference for Swagger 2.0 vs OpenAPI 3.0? → If 3.0 is required, hand-written YAML becomes the better option
