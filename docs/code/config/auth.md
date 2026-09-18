# `config/auth.go`

Reads auth-related settings from the environment.

## Exports

- `JWTSecret() []byte` — **fail-closed**: exits unless `JWT_SECRET` is 32+ bytes (the known dev value counts as unset). `ALLOW_INSECURE_JWT=true` opts into the insecure fallback for local dev only.
- `JWTTTL() time.Duration` — parses `JWT_TTL_HOURS` (default 24h); warns and falls back on invalid values.
- `AppURL() string` — public base URL for links in emails (`APP_URL`, default `http://localhost:8080`).

## Notes

- Read on every call (no caching), so tests can swap values with `t.Setenv`.
