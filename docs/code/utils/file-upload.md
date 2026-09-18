# `utils/file_upload.go`

Validation and paths for user-supplied file uploads.

## Constants & errors

- `MaxUploadSize = 5 MB`, `UploadBaseDir = "uploads"` (served publicly at `/uploads`), `AvatarUsersDir = "avatars/users"`, `AvatarContactsDir = "avatars/contacts"`.
- `ErrFileTooLarge`, `ErrInvalidFileExt` sentinels.

## Functions

- `ValidateUploadSize(size)` — rejects over-limit files.
- `ValidateAvatarExt(filename)` — allowlist `.jpg .jpeg .png .gif .webp` (case-insensitive).
- `ValidateImportExt(filename, allowed)` — e.g. `.csv` / `.json` per endpoint.
- `AvatarTarget(subdir, originalName) (urlPath, diskPath, err)` — UUID filename (no collisions/traversal), `MkdirAll`, returns the public URL path and absolute disk path.
