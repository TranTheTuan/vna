# Phase 04 — Tests & Verification

**Context:** `plans/260428-health-check-endpoints/plan.md`
**Priority:** P2 | **Status:** Complete | **Effort:** 0.1h
**Blocked by:** Phase 01, Phase 02, Phase 03

## Overview

Manual and automated verification of health endpoints.

## Requirements

- `GET /healthz` always returns 200 with `{"status":"ok"}`
- `GET /readyz` returns 200 when DB is up, 503 when DB is down
- Both endpoints accessible without `Authorization` header
- K8s probe-compatible response format

## Related Code Files

| Action | File |
|--------|------|
| Explore | `internal/delivery/http/health.go` |
| Create (optional) | `internal/delivery/http/health_test.go` |

## Implementation Steps

### 1. Manual Verification

Start the server:

```bash
make run
```

Test healthz (should always return 200):

```bash
curl -s http://localhost:8080/healthz
# Expected: {"status":"ok"} with HTTP 200
```

Test readyz (should return 200 when DB is up):

```bash
curl -s http://localhost:8080/readyz
# Expected: {"status":"ok"} with HTTP 200
```

Test without auth header (should still work):

```bash
curl -s -H "Authorization: Bearer faketoken" http://localhost:8080/healthz
# Expected: {"status":"ok"} with HTTP 200 (ignores auth)
```

### 2. Automated Test (optional, future)

Create `internal/delivery/http/health_test.go` with mock DB pool for readyz failure scenario.

### 3. DB Down Test (optional)

Stop PostgreSQL, verify `/readyz` returns 503:

```bash
sudo systemctl stop postgresql
curl -s -o /dev/null -w "%{http_code}" http://localhost:8080/readyz
# Expected: 503
sudo systemctl start postgresql
```

## Todo

- [ ] Start server with `make run`
- [ ] Test `GET /healthz` returns 200
- [ ] Test `GET /readyz` returns 200 (DB up)
- [ ] Verify no auth required for both endpoints
- [ ] (Optional) Test `/readyz` returns 503 when DB is down

## Success Criteria

- `curl http://localhost:8080/healthz` returns 200 with `{"status":"ok"}`
- `curl http://localhost:8080/readyz` returns 200 when DB is up
- Both endpoints return JSON with `{"status":"ok"}` format
- No JWT middleware on health endpoints (verified via curl with/without auth)

## Risk Assessment

- Low risk — manual testing only, no code changes
- If DB is down, `/readyz` returns 503 as expected (correct behavior)
