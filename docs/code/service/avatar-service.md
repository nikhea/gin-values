# `service/avatar_service.go`

Links stored avatar files to owners without touching other fields.

## Functions

- `SetUserAvatar(userID, avatarURL) (*Profile, error)` — 404s on unknown users; creates the profile when the user has none, otherwise only overwrites `AvatarURL` (unlike `UpdateProfile`, which replaces all fields).
- `SetContactAvatar(contactID, avatarURL) (*Contact, error)` — 404s on unknown contacts, then sets `AvatarURL`.
