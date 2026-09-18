# `dto/contact.go`

Contact payloads, list filters, and per-type value validation.

## Types

- `CreateContactRequest{UserID, Name?, Type (oneof email/phone/address/other), Value}` and `UpdateContactRequest` (same minus immutable `UserID`).
- `ContactFilter` — embeds `Pagination`, adds `user_id`, `type`, `search` query params.
- `ValidationError{Msg}` — marks client-caused failures so handlers return 400, detected with `errors.As`.

## `ValidateContactValue(type, value)`

- `email` → must parse via `net/mail`.
- `phone` → must match `^\+?[0-9][0-9\s\-().]{5,20}$` (e.g. `+15551234567`).
- `address`/`other` → non-empty; anything else → unknown-type error.

Called by `service/contact_service.go` on create/update. Unit-tested in `tests/contact_test.go`.
