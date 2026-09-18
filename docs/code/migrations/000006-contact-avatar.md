# Migration `000006` — `avatar_url` on `contacts`

- **Up:** adds nullable `avatar_url TEXT` for contact avatar images (served from `/uploads/avatars/contacts/`).
- **Down:** drops the column.
