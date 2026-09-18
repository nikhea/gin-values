# `tests/import_upload_test.go` — imports & avatar tests

- `TestContactCSVImport` — 2 valid + 1 invalid rows → `imported=2 failed=1` with the error pinned to row 4; plus 401 (no token), 400 (bad header, wrong extension, missing `user_id`), 404 (unknown user).
- `TestContactJSONImport` — 1 valid + 1 invalid → 1/1; malformed JSON and empty array → 400.
- `TestUserAvatarUpload` — upload creates a missing profile; URL shape `/uploads/avatars/users/*.png`; file exists on disk; static `GET` serves it 200; second upload works; `.txt` → 400; unknown user → 404.
- `TestContactAvatarUpload` — URL persisted on the contact; unknown contact → 404.
- Helpers: `tinyPNG` (1x1 PNG bytes), `multipartFile` (multipart request builder), `decodeSummary`.
- Uploaded files land in `tests/uploads/` (package working dir) and are removed by the `os.RemoveAll("uploads")` cleanup in `tests/setup_test.go`.
