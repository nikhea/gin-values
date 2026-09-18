# `handlers/auth_handler.go`

Gin handlers for the public auth endpoints plus the authenticated `/me`. Each carries Swagger annotations (`@Tags auth`, request/response schemas, `@Security` only on `/me`).

## Handlers

- `Register` — `POST /auth/register`: binds `dto.RegisterRequest`, 201 with message + user (no secrets leak — model uses `json:"-"`), 409 on duplicate.
- `Login` — `POST /auth/login`: 200 with `{message, token, user}`; delegates auth errors to `writeAuthError`.
- `VerifyEmail` — `GET /auth/verify?token=`: link flow, 200 on success.
- `VerifyOTP` — `POST /auth/verify-otp`: binds `dto.VerifyOTPRequest`, 200 on success (400/429 via `writeAuthError`).
- `ResendVerification` / `ForgotPassword` — always-success messages (no account enumeration).
- `ResetPassword` — accepts the token in body or `?token=` query; 200 on success.
- `Me` — `GET /auth/me` (protected): reads the ID via `middleware.GetUserID(c)` (the `req.userId` equivalent), 404 if the user vanished.

## `writeAuthError(c, err)` (private)

Maps service sentinels to status codes: 409 taken, 401 credentials, 403 unverified, 400 token/OTP/verified, 429 too many attempts, 500 otherwise.
