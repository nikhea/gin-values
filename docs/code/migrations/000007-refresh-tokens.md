# Migration `000007` — `refresh_tokens`

- **Up:** `refresh_tokens(id PK, user_id FK → users ON DELETE CASCADE, token_hash UNIQUE NOT NULL, expires_at, revoked_at NULL, replaced_by NULL, created_at)` + indexes on `token_hash` and `user_id`.
- **Down:** drops indexes and table.
