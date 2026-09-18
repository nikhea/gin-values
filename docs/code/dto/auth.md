# `dto/auth.go`

Request/response shapes for the auth endpoints, with Gin binding rules and Swagger `example` tags.

## Types

- `RegisterRequest{Name, Email, Password (8–72 chars), Age (18–100)}` — all required.
- `LoginRequest{Email, Password}` — password min 1 char (don't reveal policy on login).
- `VerifyOTPRequest{Email, Code}` — `Code` must be 6 numeric chars (`len=6,numeric`).
- `ForgotPasswordRequest{Email}`, `ResendVerificationRequest{Email}`.
- `ResetPasswordRequest{Token, NewPassword}` — token may also arrive as `?token=` query param (handled in the handler).
- `AuthResponse{Message, Token (omitempty), User}` — login/register/verify/me envelope.

## Notes

- `example:"..."` tags pre-fill Swagger UI "Try it out" payloads; regenerate docs with `swag init` after changing them.
- bcrypt caps passwords at 72 bytes, hence `max=72`.
