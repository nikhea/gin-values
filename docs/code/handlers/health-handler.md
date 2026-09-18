# `handlers/health_handler.go`

Public liveness probe for Docker healthchecks and load balancers.

## `Health`

`GET /health` (registered directly in `routes/routes.go`, outside the auth groups) → `200 {"status": "ok"}`. No auth, no DB touch — it only proves the process serves HTTP.
