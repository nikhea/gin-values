# `repository/contact_repository.go`

Thin GORM data-access layer for contacts, including filtered pagination.

## Functions

- `CreateContact` / `GetContactByID` / `UpdateContact` / `DeleteContact` — standard CRUD by contact ID.
- `ListContacts(filter dto.ContactFilter) ([]Contact, total, error)` — optional `user_id` and `type` filters, `ILIKE` search over `name`/`value`, `COUNT` for the total, then `ORDER BY created_at DESC` with `LIMIT`/`OFFSET` from `filter.Pagination` (callers must `Normalize()` first — `service/contact_service.go` does).

## Notes

- Returns both the page and the total so handlers can build `dto.NewPageMeta`.
