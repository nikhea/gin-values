# `dto/create_user.go` / `dto/update_user.go`

Admin-style user payloads (no password — auth credentials go through `dto/auth.go`).

## Types

- `CreateUserRequest{Name (3–50), Email, Age (18–100)}` — all required.
- `UpdateUserRequest` — identical shape, used by `PUT /users/:id`.

Both carry Swagger `example` tags (`"Kaige Saif"`, `kaige@example.com`, `30`).

## Notes

- Users created via `POST /users/` have no `PasswordHash` and therefore can never log in (see `service/auth_service.go` `Login`); real signup is `POST /auth/register`.
