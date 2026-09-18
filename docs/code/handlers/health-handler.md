# `handlers/health_handler.go`

Public liveness/readiness probes for Docker healthchecks and load balancers. See `handlers/errors-probes.md` for details.

- `Health` — `GET /health`, no auth, no DB touch.
- `Readyz` — `GET /readyz`, requires Postgres ping + live River client (503 otherwise).

