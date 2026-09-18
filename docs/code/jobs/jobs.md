# `jobs/jobs.go`

Background email delivery with River (Postgres-backed job queue), so SMTP latency/outages never block HTTP handlers.

## Payload & worker

- `SendEmailArgs{To, Subject, Body, HTML}` — `Kind()` returns `"send_email"`. Self-contained: the worker needs no DB state.
- `SendEmailWorker` — `Work` calls `utils.SendMail`. Failures retry with River's backoff (up to 25 attempts); states visible in the `river_job` table.

## Lifecycle

- `Setup(ctx, databaseURL)` — dedicated pgx pool → `rivermigrate` up (River owns its `river_*` tables) → registers workers → `river.NewClient` (10 workers, default queue) → `Start`. Called in `main()` and in `tests/setup_test.go` (test DB).
- `Shutdown(ctx)` — stops the client, closes the pool; nil-safe.
- `Pool` / `Client` package globals hold the live instances.

## Enqueue helpers (used by `service/auth_service.go`)

- `EnqueueSendEmail(ctx, to, subject, body, html)` — errors if the client isn't started.
- `EnqueueVerificationEmail(ctx, to, name, token, otp)` / `EnqueuePasswordResetEmail(ctx, to, name, token)` — build content via `utils` builders, then enqueue.

## Retention

- `CompletedJobRetention = 7 days`; `CleanupArgs` (`"purge_completed_jobs"`) + `CleanupWorker` delete `completed`/`discarded` rows older than that. Registered as a daily River periodic job (`PeriodicInterval(24h)`, ID `purge-completed-jobs` — idempotent across replicas).

## Notes

- River uses its own pgx pool alongside GORM — see `config/database.md`.
- River v0.47 (`river`, `riverdriver/riverpgxv5` in `go.mod`).
