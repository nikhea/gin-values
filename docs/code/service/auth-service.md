# `service/auth_service.go`

Auth business logic: registration, email verification (link + OTP), login/JWT issuance, and password reset. Sends nothing directly — all emails are enqueued as River jobs.

## Sentinel errors

`ErrInvalidCredentials` (401), `ErrEmailTaken` (409), `ErrEmailNotVerified` (403), `ErrInvalidToken` / `ErrInvalidOTP` (400), `ErrAlreadyVerified` (400), `ErrTooManyAttempts` (429). Handlers map these to status codes in `writeAuthError`.

## Functions

- `Register(req)` — rejects duplicate email, bcrypt-hashes the password, stores a verification token + OTP hash (10-min expiry), creates the user, enqueues the verification email. Enqueue failure is logged, never fatal.
- `VerifyEmail(token)` — link flow: marks verified, clears token; `ErrAlreadyVerified` if done.
- `VerifyOTP(email, code)` — code flow: unknown email → generic `ErrInvalidOTP` (no enumeration); checks expiry, caps guesses at `utils.OTPMaxAttempts` (then `ErrTooManyAttempts` until resend); success clears token + OTP fields and verifies.
- `ResendVerification(email)` — rotates token + OTP, resets attempts, re-enqueues; unknown emails silently succeed.
- `Login(req)` — email lookup → empty-hash guard (pre-auth users can't log in) → bcrypt check → verified check → `utils.GenerateToken`.
- `ForgotPassword(email)` — sets 1-hour reset token, enqueues reset email; unknown emails silently succeed.
- `ResetPassword(token, newPassword)` — validates token + expiry, hashes the new password, clears reset fields (single-use).

## Depends on

`repository`, `utils` (password/JWT/OTP), `jobs` (email enqueue), `dto`, `models`.
