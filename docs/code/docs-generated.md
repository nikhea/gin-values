# `docs/docs.go`, `swagger.json`, `swagger.yaml`

Generated Swagger API documentation — output of `swag init`, never edited by hand.

- Source of truth: handler annotations (`@Summary`, `@Param`, `@Success`, `@Security BearerAuth`) plus `example` struct tags on DTOs/models.
- Served at runtime: `GET /swagger/*any` (wired in `routes/routes.go`).
- After changing annotations or examples, regenerate with `swag init` from the repo root.
