# `config/env.go`

Loads environment variables from the `.env` file in the repo root.

## Exports

- `LoadEnv()` — calls `godotenv.Load()`; logs `slog.Info("No .env file found")` when absent (normal in production, where real env vars are injected).

## Environment variables used by the app

| Variable | Used in | Purpose |
|---|---|---|
| `DATABASE_URL` | `config/database.go`, `config/migrate.go`, `main.go` (River) | Postgres DSN |
| `APP_PORT` | — (informational; server currently listens on `:8080`) | — |
| `APP_URL` | `config/auth.go` (`AppURL`) | Public base URL for email links |
| `JWT_SECRET` | `config/auth.go` | HMAC secret for access tokens |
| `JWT_TTL_HOURS` | `config/auth.go` | Legacy access-token lifetime (hours); overridden by `ACCESS_TTL_MINUTES` |
| `ACCESS_TTL_MINUTES` | `config/auth.go` (`AccessTTL`) | Access JWT lifetime, default `15` |
| `REFRESH_TTL_DAYS` | `config/auth.go` (`RefreshTTL`) | Refresh token lifetime, default `30` |
| `EMAIL_SERVICE` | `config/mail.go` | `Gmail` → `smtp.gmail.com:587` |
| `EMAIL_ADDRESS` / `EMAIL_PASSWORD` | `config/mail.go` | SMTP credentials (app password) |
| `EMAIL_HOST` / `EMAIL_PORT` | `config/mail.go` | Override when service isn't Gmail |
| `TEST_DATABASE_URL` | `tests/setup_test.go` | Test database DSN |

## Notes

- `.env` is gitignored — secrets never enter version control. Values are only listed here by name, never by content.
