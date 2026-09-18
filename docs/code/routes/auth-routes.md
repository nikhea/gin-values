# `routes/auth_routes.go`

Public auth route table plus the self-guarded `/me`.

## `RegisterAuthRoutes(router)`

Mounts on `/api/auth`: `POST register`, `POST login`, `GET verify`, `POST verify-otp`, `POST resend-verification`, `POST forgot-password`, `POST reset-password`, and `GET me` with inline `middleware.AuthRequired()`.

`GET /me` is registered here (next to the other auth endpoints) rather than in the protected group, guarding itself so the auth routes file stays self-contained.
