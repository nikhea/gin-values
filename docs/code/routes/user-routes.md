# `routes/user_routes.go`

CRUD route table for users, mounted on the protected `/users` group.

## `RegisterUserRoutes(api)`

`POST /`, `GET /`, `GET /:id`, `PUT /:id`, `DELETE /:id` → `handlers.CreateUser`, `GetUsers`, `GetUser`, `UpdateUser`, `DeleteUser`. Auth comes from the parent group in `routes/routes.go`, so nothing is added here.
