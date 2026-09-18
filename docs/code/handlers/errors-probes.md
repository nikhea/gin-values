# `handlers/errors.go` + `handlers/health_handler.go` (probes)

## `internalError(c, err)` (`errors.go`)

Single funnel for all 500s: logs the real error with `requestID`, method and route via `slog.Error`, responds `{"error": "internal server error"}`. Every handler 500 (including `writeAuthError`'s default branch) routes through it — SQL/driver/filesystem details never reach clients.

## Probes (`health_handler.go`)

- `Health` — `GET /health`: liveness only, no auth, no DB touch (Docker `HEALTHCHECK`).
- `Readyz` — `GET /readyz`: pings Postgres (`SELECT`-level via `PingContext`) and requires a live River client; `503 {"error": ...}` naming the unready dependency otherwise. Point load-balancer checks here.
