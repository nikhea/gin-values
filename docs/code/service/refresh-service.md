# `models/refresh_token.go` + `repository/refresh_repository.go` + `service/refresh_service.go`

Opaque rotating refresh tokens (Postgres-backed sessions).

## Model

`RefreshToken{ID, UserID (cascade), TokenHash (unique), ExpiresAt, RevokedAt?, ReplacedBy?, CreatedAt}` — everything `json:"-"`. Only SHA-256 hashes are stored; the plaintext travels once (login/refresh response) and is never persisted.

## Repository

`CreateRefreshToken`, `GetRefreshTokenByHash`, `GetRefreshTokenByID`, `UpdateRefreshToken`, `RevokeUserTokens(userID)` (logout-all, one UPDATE), `RevokeTokenFamily(start)` (walks `ReplacedBy` links revoking the chain).

## Service

- `IssueTokenPair(user)` — fresh access JWT + single-use refresh token.
- `Rotate(raw)` — validates (unknown/expired → `ErrInvalidRefreshToken`), loads the user, issues a new pair, revokes the old with `ReplacedBy` set. A presented token that was already rotated signals theft **only while its replacement is still alive** → family revoke + `ErrRefreshReuse`; otherwise plain invalid (covers logout/logout-all without false alarms).
- `Revoke(raw)` — logout, silent on unknown tokens (no validity oracle).
- `RevokeAll(userID)` — logout everywhere.

Sentinels: `ErrInvalidRefreshToken` / `ErrRefreshReuse` → 401 in `writeAuthError`. Endpoints: `POST /auth/refresh`, `/auth/logout` (public, rate-limited), `/auth/logout-all` (protected). TTLs: `config.AccessTTL()` (15m) / `config.RefreshTTL()` (30d). Tested in `tests/refresh_test.go` (rotation, reuse-kills-family, expiry, logout semantics, HTTP).
