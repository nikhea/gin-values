# `models/users.go`

GORM model for the `users` table, including auth columns.

## `User` fields

| Field | Column / notes |
|---|---|
| `ID` | `varchar(36)` PK, UUID string |
| `Name`, `Email` (unique), `Age` | Basic profile data |
| `PasswordHash` | `json:"-"` — never serialized to API responses |
| `EmailVerified` | `bool` — gates login until verification |
| `VerificationToken` | `json:"-"` — email link token (cleared on verify) |
| `ResetToken`, `ResetExpiresAt` | `json:"-"` — 1-hour password-reset token |
| `OTPHash`, `OTPExpiresAt`, `OTPAttempts` | `json:"-"` — SHA-256 of the 6-digit code, expiry, guess counter |
| `CreatedAt`, `UpdatedAt` | Timestamps |
| `Profile *Profile` | One-to-one, `OnDelete:CASCADE`, `omitempty` |

## Notes

- Every secret column uses `json:"-"` so hashes/tokens can't leak through handlers that return the model directly.
- Schema comes from migrations `000001`, `000004`, `000005` (see `migrations/` docs).
