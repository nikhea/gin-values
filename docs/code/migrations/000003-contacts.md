# Migration `000003` — create `contacts`

- **Up:** `contacts(id PK, user_id NOT NULL, name, type NOT NULL, value NOT NULL, timestamps)` with FK `fk_contacts_user → users(id) ON DELETE CASCADE`, plus indexes on `user_id` and `type` (both used by `ListContacts` filters).
- **Down:** drops indexes, constraint, and table.
