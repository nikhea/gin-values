# `handlers/contact_handler.go`

Gin handlers for contacts, including the filtered/paginated list. All JWT-protected.

## Handlers

- `CreateContact` — `POST /contacts/`: 201; `writeContactError` maps validation errors to 400 and missing users to 404.
- `GetContacts` — `GET /contacts/`: binds `dto.ContactFilter` query (`user_id`, `type`, `search`, `page` default 1, `page_size` default 10 — all with Swagger `example`s), returns `{message, contacts, meta}`.
- `GetContact` / `UpdateContact` / `DeleteContact` — by `:id`, standard 200/400/404/500 mapping.

## `writeContactError(c, err, notFoundMsg)` (private)

`errors.As` → `*dto.ValidationError` gives 400; `errors.Is` → `gorm.ErrRecordNotFound` gives 404; else 500.
