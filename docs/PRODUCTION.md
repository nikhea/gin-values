# Production runbook

Checklist and operating notes for running Grip in production. Start from `.env.example` — every variable is documented there.

## 1. Pre-deploy checklist

- [ ] `JWT_SECRET` is 32+ random bytes (`openssl rand -hex 32`). The app **refuses to boot** without it (fail-closed); `ALLOW_INSECURE_JWT=true` is dev-only.
- [ ] `DATABASE_URL` points at managed Postgres. If you sit a pooler (PgBouncer) in front, prefer **session** mode — GORM/pgx prepared statements and River's `LISTEN/NOTIFY` both misbehave behind transaction-mode pooling. Give River a direct connection if in doubt.
- [ ] `APP_URL` is the public base URL (email links depend on it).
- [ ] `TRUSTED_PROXIES` lists your LB/CDN ranges so rate limiting and logs see real client IPs. Empty = trust none (safe default for direct exposure).
- [ ] `CORS_ALLOWED_ORIGINS` lists exactly your frontend origin(s). Empty = no CORS headers.
- [ ] `GIN_MODE=release` (set automatically in the prod image).
- [ ] `EMAIL_*` set (Gmail app password, not your login password), else emails only log.

## 2. TLS / reverse proxy

Terminate TLS at the proxy; the app serves plain HTTP. Minimal Caddy example:

```caddy
api.example.com {
    reverse_proxy 127.0.0.1:8080
}
```

Forward `X-Forwarded-For` (Caddy/nginx do by default) and add the proxy's range to `TRUSTED_PROXIES`.

## 3. Migrations — run once, not per replica

The server runs `golang-migrate` + River migrations on boot. With **one** replica that's fine. With 2+ replicas, run migrations as a one-shot step instead and keep boot migration as a safety net:

```bash
# Docker: the image ships a migrate binary
docker compose run --rm api /app/migrate up
# or: go run ./cmd/migrate up
```

## 4. Backups

Postgres holds everything (app data + River jobs). Nightly `pg_dump` at minimum:

```bash
pg_dump "$DATABASE_URL" -Fc -f "gin-backup-$(date +%F).dump"
# restore: pg_restore -d "$DATABASE_URL" gin-backup-<date>.dump
```

Back up the `uploads/` volume on the same schedule (avatars aren't in the DB).

## 5. Background jobs (River)

- Each API replica runs its own River client/workers — this is supported and how jobs scale horizontally.
- River requires a **direct** Postgres connection (no transaction-mode pooler) for `LISTEN/NOTIFY`; without it, workers still function via polling.
- Finished jobs are purged automatically: a daily River periodic job (`purge_completed_jobs`) deletes `completed`/`discarded` rows older than 7 days (`jobs.CompletedJobRetention`).
- Dead (`discarded`) jobs younger than that are visible in `river_job.errors` — alert on growth.

## 6. Observability

- Logs are JSON on stdout (`slog`) — ship them to your log aggregator; `requestID` joins access and error lines.
- Probes: `/health` (liveness — process is up) vs `/readyz` (readiness — DB + River reachable). Point load-balancer health checks at `/readyz`, container `HEALTHCHECK` at `/health`.
- Alert on: 5xx rate, 429 rate spikes (possible abuse or too-tight budgets), `river_job` growth, login 401 spikes.

## 7. Security recap (enforced in code)

- Auth endpoints: 10 req/min/IP (`RATE_LIMIT_AUTH_PER_MINUTE`); API: 300 req/min/IP (`RATE_LIMIT_API_PER_MINUTE`).
- 500s return a generic message; details go to logs with `requestID`.
- Server timeouts: 10s read/header, 30s write, 60s idle. Uploads capped at 5 MB with an image allowlist.
- OTP: 6 digits, 10-min expiry, SHA-256 stored, 5-guess lockout. Passwords: bcrypt, 8–72 chars.
