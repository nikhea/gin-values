# `models/contact.go`

GORM model for the `contacts` table, plus the allowed contact types.

## Constants

`ContactTypeEmail = "email"`, `ContactTypePhone = "phone"`, `ContactTypeAddress = "address"`, `ContactTypeOther = "other"` — mirrored by the `oneof=` binding tags in `dto/contact.go`.

## `Contact` fields

`ID` (UUID PK), `UserID` (indexed FK → `users.id`), `Name`, `Type` (indexed, `varchar(16)`), `Value`, `AvatarURL` (`avatar_url`, `omitempty` — public `/uploads/...` path set by the avatar endpoint), `CreatedAt`, `UpdatedAt`.

## Notes

- Filtered listing (`user_id`, `type`, ILIKE search) is implemented in `repository/contact_repository.go` (`ListContacts`); value validation per type lives in `dto/contact.go` (`ValidateContactValue`).
