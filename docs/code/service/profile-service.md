# `service/profile_service.go`

One-to-one profile business logic, addressed by user ID.

## Functions

- `CreateProfile(userID, req)` — 404s (via `ErrRecordNotFound`) when the user doesn't exist, then inserts with a fresh UUID.
- `GetProfileByUserID`, `UpdateProfile` (overwrite all fields + save), `DeleteProfile` (existence check first for 404s).

## Depends on

`repository`, `dto`, `models`, `google/uuid`.
