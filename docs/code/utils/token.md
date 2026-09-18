# `utils/token.go`

Cryptographically random hex tokens for email-verification and password-reset links.

## `GenerateSecureToken(n) (string, error)`

Reads `n` bytes from `crypto/rand` and hex-encodes them (32 bytes → 64-char token). Used by `service/auth_service.go` for `VerificationToken`/`ResetToken`. Unit-tested in `tests/token_test.go` (length, hex validity, uniqueness).
