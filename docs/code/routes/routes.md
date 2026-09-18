# `routes/routes.go`

Assembles the Gin engine: proxy trust, middleware, Swagger UI, probes, public auth routes, and the rate-limited + JWT-protected API groups.

## `Setup() *gin.Engine`

1. `gin.Default()` + `SetTrustedProxies(trustedProxies())` — `TRUSTED_PROXIES` env (IPs/CIDRs); empty trusts none so `ClientIP` falls back to the peer. Misconfigured values are fatal.
2. `middleware.RequestID()` + `corsMiddleware()` — CORS emits headers only for explicitly configured `CORS_ALLOWED_ORIGINS` (methods GET/POST/PUT/DELETE/OPTIONS, `Authorization` allowed, credentials on, 12h max-age).
3. `GET /swagger/*any`, public `GET /health` + `GET /readyz`, static `/uploads`.
4. `RegisterAuthRoutes(router)` — with the strict `"auth"` rate-limit scope.
5. Protected `/api` group: `"api"` rate limit first, then `AuthRequired()`, then users/profiles/contacts.

## Helpers

`trustedProxies()` / `corsMiddleware()` parse their env vars (comma-separated, trimmed); both default to the closed posture when unset.
