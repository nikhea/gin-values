# `repository/user_repository.go`

Thin GORM data-access layer for users. No business logic — just queries against `config.DB`.

## Functions

- `CreateUser(user)` — insert.
- `GetUsers()` — all users with `Profile` preloaded.
- `GetUserByID(id)` — one user with `Profile` preloaded; `gorm.ErrRecordNotFound` when missing.
- `GetUserByEmail(email)` — lookup for login/register/OTP flows (with `Profile`).
- `GetUserByVerificationToken(token)` / `GetUserByResetToken(token)` — token-link flows.
- `UpdateUser(user)` — full-row `Save`.
- `DeleteUser(id)` — delete by ID (profile cascades at the DB level).

## Notes

- Services translate `ErrRecordNotFound` into 404s or auth sentinel errors; handlers never touch `config.DB` directly.
