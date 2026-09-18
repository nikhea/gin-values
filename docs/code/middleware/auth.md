# `middleware/auth.go`

JWT gate for protected routes, plus the `req.userId`-style accessors.

## `AuthRequired() gin.HandlerFunc`

Reads the `Authorization` header, accepting both `Bearer <token>` and a bare `<token>` (wrong schemes rejected). Validates signature + expiry via `utils.ParseToken`; on failure aborts 401. On success stores the JWT subject/email in the Gin context and calls `Next()`.

## Accessors (use inside any protected handler)

- `GetUserID(c) (string, bool)` — the authenticated user's ID (`claims.Subject`).
- `GetUserEmail(c) (string, bool)` — the authenticated user's email.
- `ContextUserID` / `ContextUserEmail` — raw context keys (prefer the accessors).

`ok == false` only happens if the middleware was bypassed. Unit-tested in `tests/middleware_test.go`; wired onto the `/api` group in `routes/routes.go`.
