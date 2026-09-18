# `middleware/rate_limit.go`

Fixed-window per-IP rate limiting (stdlib only: `sync.Map` + mutex, no new dependency).

## Budgets

`AuthRateLimit = 10` req/min (anonymous auth endpoints — brute-force protection), `APIRateLimit = 300` req/min (authenticated API), shared `RateWindow = 1 minute`. Overridable per deploy/test via `RATE_LIMIT_<SCOPE>_PER_MINUTE` (e.g. `RATE_LIMIT_AUTH_PER_MINUTE`; hyphens in scopes normalize to underscores).

## `RateLimit(scope, n)`

Keys buckets by `scope + ClientIP` (correct behind proxies only when `TRUSTED_PROXIES` is set — see `routes/routes.md`). Consumes one token per request; over budget → `429 + Retry-After: 60` with a generic message. A `sync.Once` janitor goroutine purges expired buckets every 5 minutes so rotating IPs can't grow the map.

Wired in `routes/auth_routes.go` (`"auth"` scope) and on the protected `/api` group *before* `AuthRequired` (cheap rejection first). Tested in `tests/rate_limit_test.go`; the integration suite disables limits via env in `tests/setup_test.go`.
