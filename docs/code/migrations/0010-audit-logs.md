# Migration `0010` — `audit_logs`

- **Up:** `audit_logs(id PK, actor_id NULL, action, resource_type, resource_id NULL, org_id NULL, ip, request_id NULL, metadata JSONB '{}', created_at)` + indexes on `actor_id`, `org_id`, `action`, `created_at`.
- **Down:** drops indexes and table. (Rows are append-only by convention — the app has no delete path.)
