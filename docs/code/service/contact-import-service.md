# `service/contact_import_service.go`

Bulk contact creation from CSV and JSON streams.

## Functions

- `ImportContactsFromCSV(userID, r)` — requires `name,type,value` header (case-insensitive); data rows start at line 2; wrong column counts become per-row errors (`expected 3 columns...`).
- `ImportContactsFromJSON(userID, r)` — decodes a `[{name,type,value}]` array; empty array is an error.
- `importRows(...)` (private) — validates each row with `dto.ValidateContactValue`, creates valid contacts with UUIDs, collects `{row, message}` failures into `dto.ImportSummary{Imported, Failed, Errors}`.

Both verify the user exists first (`ErrRecordNotFound` → 404 in handlers). Covered by `tests/import_upload_test.go`.
