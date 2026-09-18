# `handlers/contact_import_handler.go` + `handlers/avatar_handler.go`

Multipart upload handlers (all JWT-protected).

## Import (`contact_import_handler.go`)

- `POST /contacts/import/csv` (`ImportContactsCSV`) and `POST /contacts/import/json` (`ImportContactsJSON`), sharing the private `importContacts` helper: `user_id` from form field or query (400 if missing), `file` multipart field (400 if missing), 5 MB + extension checks, then the matching service parser. Unknown user → 404; unparsable file → 400; partial success → 200 `dto.ImportSummary`.

## Avatars (`avatar_handler.go`)

- `POST /users/:id/avatar` (`UploadUserAvatar`) and `POST /contacts/:id/avatar` (`UploadContactAvatar`), sharing private `uploadAvatar`: validates size + image extension, stores via `utils.AvatarTarget` + `c.SaveUploadedFile`, links through `service.SetUserAvatar`/`SetContactAvatar` (404 on unknown owner, orphan file removed on link failure), returns `dto.AvatarResponse`.
