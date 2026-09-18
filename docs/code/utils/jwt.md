# `utils/jwt.go`

JWT issuance and validation (HS256, `golang-jwt/jwt/v5`).

## Types & functions

- `Claims{Email, jwt.RegisteredClaims}` — user ID in `Subject`, email as a convenience claim.
- `GenerateToken(userID, email) (string, error)` — `IssuedAt` now, `ExpiresAt` now + `config.JWTTTL()`, signed with `config.JWTSecret()`.
- `ParseToken(tokenString) (*Claims, error)` — rejects unexpected signing methods, bad signatures, expired/malformed tokens, and empty subjects.

Used by `service/auth_service.go` (`Login`) and `middleware/auth.go` (`AuthRequired`). Unit-tested in `tests/jwt_test.go` (round-trip, tampered, wrong secret, expired).
