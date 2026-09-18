# `config/auth.go`

Reads auth-related settings from the environment.

## Exports

- `JWTSecret() []byte` — returns `JWT_SECRET`. Falls back to an insecure dev-only constant with a `slog.Warn` when unset (keeps local dev working; never rely on it in production).
- `JWTTTL() time.Duration` — parses `JWT_TTL_HOURS` (default 24h); warns and falls back on invalid values.
- `AppURL() string` — public base URL for links in emails (`APP_URL`, default `http://localhost:8080`).

## Notes

- Read on every call (no caching), so tests can swap values with `t.Setenv`.
