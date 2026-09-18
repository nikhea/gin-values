# `models/profile.go`

GORM model for the `profiles` table — one-to-one extension of a user.

## `Profile` fields

`ID` (UUID PK), `UserID` (unique, non-null FK → `users.id`), `Bio`, `AvatarURL`, `Phone`, `Address`, `CreatedAt`, `UpdatedAt`.

## Notes

- Uniqueness on `user_id` enforces at most one profile per user.
- Deleting a user cascades to its profile (FK `ON DELETE CASCADE`, migration `000002`).
