# `dto/import.go`

Shapes for bulk contact imports and avatar uploads.

## Types

- `ContactImportRow{Name, Type (oneof email/phone/address/other), Value}` — one CSV row or JSON array element, with Swagger examples.
- `RowError{Row, Message}` — one skipped row (1-based source line).
- `ImportSummary{Message, Imported, Failed, Errors?}` — valid rows are created, invalid ones skipped and listed.
- `AvatarResponse{Message, AvatarURL}` — returned by avatar endpoints with the public `/uploads/...` URL.
