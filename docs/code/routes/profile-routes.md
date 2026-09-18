# `routes/profile_routes.go`

One-to-one profile routes nested under users: `POST /:id/profile`, `GET /:id/profile`, `PUT /:id/profile`, `DELETE /:id/profile` → the matching `handlers` profile functions. Takes the protected `/users` group; `:id` is the user ID.
