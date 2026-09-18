# `routes/routes.go`

Assembles the Gin engine: global middleware, Swagger UI, public auth routes, and the JWT-protected API groups.

## `Setup() *gin.Engine`

1. `gin.Default()` + `middleware.RequestID()`.
2. `GET /swagger/*any` — Swagger UI.
3. `RegisterAuthRoutes(router)` — public `/api/auth/*` (plus `/me`, which guards itself).
4. Protected group on `/api` with `middleware.AuthRequired()`:
   - `/users` → `RegisterUserRoutes` + `RegisterProfileRoutes` (nested `/:id/profile`)
   - `/contacts` → `RegisterContactRoutes` (takes the protected group, so contacts inherit auth)

Every non-auth endpoint therefore requires `Authorization: Bearer <token>` (or bare token).
