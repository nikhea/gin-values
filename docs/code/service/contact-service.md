# `service/contact_service.go`

Contact business logic: validation, ownership checks, and paged listing.

## Functions

- `CreateContact(req)` — validates `type`/`value` (`dto.ValidateContactValue`), ensures the user exists (else 404), inserts with UUID.
- `ListContacts(filter)` — normalizes pagination, delegates to `repository.ListContacts`.
- `GetContactByID`, `UpdateContact` (re-validate + overwrite + save), `DeleteContact` (existence check first for 404s).

## Depends on

`repository`, `dto`, `models`, `google/uuid`.
