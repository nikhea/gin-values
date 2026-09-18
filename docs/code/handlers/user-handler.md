# `handlers/user_handler.go`

Gin handlers for user CRUD. All routes are JWT-protected (`@Security BearerAuth` on every operation).

## Handlers

- `CreateUser` — `POST /users/`: binds `dto.CreateUserRequest`, 201. Creates password-less users (see `dto/user-dtos.md`).
- `GetUsers` — `GET /users/`: 200 list, 404 `{message: "No users yet"}` when empty.
- `GetUser` — `GET /users/:id`: 200 user, 404 when missing.
- `UpdateUser` — `PUT /users/:id`: binds + updates, 400/404/500.
- `DeleteUser` — `DELETE /users/:id`: 200 confirmation, 404 when missing.

## Conventions

`errors.Is(err, gorm.ErrRecordNotFound)` → 404; bind failures → 400 with the error; everything else → 500. Response shapes reference `dto/responses.md` envelopes.
