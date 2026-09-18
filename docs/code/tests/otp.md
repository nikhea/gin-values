# `tests/otp_test.go` — OTP integration tests

- `TestVerifyOTPFlow` — register stores hash + future expiry; wrong code → `ErrInvalidOTP` with attempt counted; unknown email indistinguishable; correct code verifies → login works; re-verify → `ErrAlreadyVerified`.
- `TestVerifyOTPExpiryAndLockout` — expired code fails; 5 wrong guesses → `ErrTooManyAttempts` (correct code also locked); resend rotates + resets, then verifies.
- `TestVerifyOTPHTTP` — `POST /auth/verify-otp` returns 400/200 through `routes.Setup()`.
- `setKnownOTP` helper plants a deterministic code (`HashOTP("123456")`) since real codes only travel by email.
