# Grip

A REST API built with **Go + Gin + Postgres**, covering users, profiles, and contacts with JWT auth, email verification (link + OTP), background email jobs, file uploads, and Swagger docs.

## Features

- **Users / Profiles / Contacts** — full CRUD, one-to-one profiles, filtered + paginated contact listing
- **Auth** — register, login (JWT `Bearer`), email verification via link **or** 6-digit OTP (10-min expiry, 5-attempt lockout), password reset, resend flows
- **Email jobs** — River (Postgres-backed queue) sends verification/reset emails async; HTML + plain-text templates
- **Cache** — optional Redis layer (graceful Postgres fallback) with namespaced keys; hot user reads cached, invalidated on writes
- **Cache** — optional Redis layer (graceful Postgres fallback) with namespaced keys; user reads cached with invalidation on writes
- **Imports** — bulk contact import from CSV (`name,type,value`) or a JSON array, with per-row error reports
- **Avatars** — image uploads for users and contacts, served from `/uploads`
- **Docs** — Swagger UI at `/swagger/*any`, per-file guides in `docs/code/`
- **Ops** — structured JSON logs (`slog`), `/health` probe, multi-target Docker builds, GitHub Actions CI

## Tech stack

Go 1.27 · Gin · GORM · Postgres · Redis · River · golang-migrate · JWT (HS256) · bcrypt · Swagger · Docker

## Prerequisites

- Go 1.27+
- Postgres 16 (local) **or** Docker
- Optional: `swag`, `golangci-lint`, `air` (all installable via `go install`; see `Makefile`)

## Quickstart (local)

```bash
# 1. Create .env (DATABASE_URL, JWT_SECRET, EMAIL_* — see config/env.go docs)

# 2. Create database + run migrations + start
createdb gin_app
go run ./cmd/migrate up
go run .
```

API listens on `:8080`. Swagger UI: http://localhost:8080/swagger/index.html

Key env vars:

| Variable | Purpose | Default |
|---|---|---|
| `DATABASE_URL` | Postgres DSN | — (required) |
| `APP_PORT` | HTTP listen port | `8080` |
| `JWT_SECRET` | Token signing secret, 32+ bytes | — (required; boot fails without it) |
| `ALLOW_INSECURE_JWT` | Permit dev fallback secret | unset (local dev only, never prod) |
| `JWT_TTL_HOURS` | Token lifetime | `24` |
| `APP_URL` | Public base URL for email links | `http://localhost:8080` |
| `EMAIL_SERVICE` / `EMAIL_ADDRESS` / `EMAIL_PASSWORD` | Gmail SMTP (app password) | unset → emails are logged, not sent |
| `TRUSTED_PROXIES` | LB/CDN IPs or CIDRs (comma-separated) | trust none |
| `CORS_ALLOWED_ORIGINS` | Browser origins (comma-separated) | none (no CORS headers) |
| `RATE_LIMIT_AUTH_PER_MINUTE` / `RATE_LIMIT_API_PER_MINUTE` | Per-IP budgets | `10` / `300` |
| `REDIS_ADDR` / `REDIS_PASSWORD` / `REDIS_DB` | Redis cache (optional — app runs without it) | `localhost:6379` / empty / `0` |
| `TEST_DATABASE_URL` | DB used by `tests/` | local `gin_app_test` |

A commented template lives in `.env.example`. For deployments, read `docs/PRODUCTION.md` (secrets, TLS, migrations, backups, jobs, observability).

## Quickstart (Docker)

```bash
# Development (live reload via Air)
docker compose -f docker-compose.yml -f docker-compose.dev.yml up --build

# Production
docker compose -f docker-compose.yml -f docker-compose.prod.yml up -d --build
```

This boots Postgres + API with healthchecks. `JWT_SECRET` should be overridden in prod (`openssl rand -hex 32`). See `docs/code/docker-compose.md`.

## API overview

Base path `/api`. Everything except `/api/auth/*` (and `/health`) needs `Authorization: Bearer <token>`.

| Method & path | Description |
|---|---|
| `POST /auth/register` | Sign up (sends verification email with OTP + link) |
| `POST /auth/login` | Log in, returns JWT |
| `GET /auth/verify?token=` | Verify email via link |
| `POST /auth/verify-otp` | Verify email with `{email, code}` |
| `POST /auth/resend-verification` | Resend verification email |
| `POST /auth/forgot-password` / `POST /auth/reset-password` | Password reset flow |
| `GET /auth/me` | Current user (protected) |
| `GET/POST /users/` · `GET/PUT/DELETE /users/:id` | User CRUD |
| `POST /users/:id/avatar` | Upload user avatar |
| `GET/POST/PUT/DELETE /users/:id/profile` | Profile endpoints |
| `GET/POST /contacts/` · `GET/PUT/DELETE /contacts/:id` | Contacts (+ `?user_id=&type=&search=&page=&page_size=`) |
| `POST /contacts/import/csv` · `POST /contacts/import/json` | Bulk import (multipart `file` + `user_id`) |
| `POST /contacts/:id/avatar` | Upload contact avatar |
| `GET /health` | Liveness probe |

## Development

```bash
make fmt      # format code
make vet      # go vet
make lint     # golangci-lint v2
make test     # full suite (needs Postgres; skips DB tests if unreachable)
make check    # fmt-check + vet + lint + test (same as CI)
make swagger  # regenerate Swagger docs after annotation changes
swag init     # same thing directly
```

Migrations live in `migrations/` (`000001`–`000006`); River manages its own `river_*` tables automatically. Roll back with `go run ./cmd/migrate down`.

## Project structure

```
main.go               # entrypoint: logger → env → migrations → DB → River → Gin
config/               # env, logger (slog JSON), database (GORM), migrations, auth/mail settings
models/               # User, Profile, Contact (secrets use json:"-")
dto/                  # request/response shapes + validation
repository/           # thin GORM queries (no business logic)
service/              # business logic (auth, users, profiles, contacts, imports, avatars)
handlers/             # Gin handlers + Swagger annotations
middleware/           # AuthRequired (+ GetUserID/GetUserEmail), RequestID
routes/               # route tables; /api/* protected except /api/auth/*
jobs/                 # River client, workers, email enqueue helpers
utils/                # JWT, bcrypt, OTP, tokens, mailer, HTML templates, uploads
migrations/           # versioned SQL (golang-migrate)
tests/                # all tests: unit + DB integration (auto-skips without Postgres)
docs/                 # generated Swagger + docs/code (one guide per source file)
```

## CI

`.github/workflows/`: `ci.yml` (fmt, vet, lint, tests with Postgres service) and `docker.yml` (builds both image targets, boots prod stack, health + register smoke test).
