# Migration `0011` — `notifications`

- **Up:** `notifications(id PK, user_id FK → users CASCADE, type, title, body, data JSONB '{}', read_at NULL, created_at)` + indexes on `user_id` and `read_at`.
- **Down:** drops indexes and table.
