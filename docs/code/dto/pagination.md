# `dto/pagination.go`

Shared page/page_size handling and the list-envelope metadata.

## Types & functions

- `Pagination{Page, PageSize}` (`form` tags) — `Normalize()` defaults to page 1 / 10 per page, clamps size to 100; `Offset()` returns `(Page-1)*PageSize` for SQL.
- `PageMeta{Page, PageSize, Total, TotalPages}` — returned as `meta` in list responses.
- `NewPageMeta(p, total)` — builds the envelope (rounds total pages up).

Used by the contacts list endpoint (`handlers/contact_handler.go` → `dto.NewPageMeta`).
