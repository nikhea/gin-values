# `dto/profile.go`

Profile payloads for create/update (all fields optional except via binding limits).

## Types

- `CreateProfileRequest` / `UpdateProfileRequest` — `Bio` (max 1000), `AvatarURL` (optional valid URL, max 2048), `Phone` (max 32), `Address` (max 500), each with Swagger `example` tags.

## Notes

- `UserID` comes from the URL (`/users/:id/profile`), never from the body — enforced by `routes/profile_routes.go` + `service/profile_service.go`.
