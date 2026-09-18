# `utils/otp.go`

One-time passcodes for email verification.

## Constants

`OTPLength = 6` digits, `OTPExpiry = 10 minutes`, `OTPMaxAttempts = 5` guesses per code.

## Functions

- `GenerateOTP() (string, error)` — random numeric code via `crypto/rand` (uniform per digit).
- `HashOTP(code) string` — SHA-256 hex digest; only this is stored (`users.otp_hash`), the plain code travels by email only.
- `CheckOTP(code, hash) bool` — constant-time comparison; empty inputs always fail.

Enforced by `service/auth_service.go` (`VerifyOTP`), which additionally tracks expiry + attempts. Unit-tested in `tests/otp_utils_test.go`.
