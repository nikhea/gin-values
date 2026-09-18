# `handlers/profile_handler.go`

Gin handlers for the nested one-to-one profile endpoints (`/users/:id/profile`). All JWT-protected.

## Handlers

`CreateProfile` (201), `GetProfile` (200), `UpdateProfile` (200), `DeleteProfile` (200 confirmation) — each takes `:id` as the **user** ID, binds `dto.Create/UpdateProfileRequest` where applicable, and maps `gorm.ErrRecordNotFound` to 404 (`"User/Profile not found"`).

A `var _ = models.Profile{}` line keeps the model visible to Swagger generation.
