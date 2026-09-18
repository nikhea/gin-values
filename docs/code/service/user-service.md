# `service/user_service.go`

User CRUD business logic (admin-style; no passwords involved).

## Functions

- `CreateUser(req dto.CreateUserRequest)` — assigns a UUID and inserts. Note: no `PasswordHash`, so such users can never log in — real signup is `service/auth_service.go` `Register`.
- `GetUserByID(id)` — cache-aside via `utils.UserKey(id)` (5-min TTL): Redis hit returns immediately, miss/error falls back to Postgres and repopulates. Cache failures are debug-logged, never fatal.
- `UpdateUser` / `DeleteUser` — invalidate `utils.UserKey(id)` after the write (TTL bounds any residual staleness). Auth mutations (`VerifyEmail`, `VerifyOTP`, `ResetPassword`) invalidate the same key so verification state is never stale.
- `UpdateUser(id, req)` — loads (404 if missing), overwrites name/email/age, saves.
- `DeleteUser(id)` — verifies existence first so callers can return 404, then deletes.

## Depends on

`repository`, `dto`, `models`, `google/uuid`. Covered by `tests/users_test.go`.
