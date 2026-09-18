# `repository/profile_repository.go`

Thin GORM data-access layer for profiles, keyed by `user_id` (not profile ID).

## Functions

- `CreateProfile(profile)` — insert.
- `GetProfileByUserID(userID)` — `ErrRecordNotFound` when the user has no profile.
- `UpdateProfile(profile)` — full-row `Save`.
- `DeleteProfileByUserID(userID)` — delete by `user_id`.

## Notes

- Services check user/profile existence first so handlers can return 404 (`service/profile_service.go`).
