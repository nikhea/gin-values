# Migration `000002` — create `profiles`

- **Up:** `profiles(id PK, user_id NOT NULL, bio, avatar_url, phone, address, timestamps)` with FK `fk_profiles_user → users(id) ON DELETE CASCADE` + unique index on `user_id` (one profile per user).
- **Down:** drops indexes, constraint, and table.
