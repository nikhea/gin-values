# Migration `000008` — soft delete

- **Up:** adds `deleted_at TIMESTAMPTZ` (+ index) to `users`, `profiles`, `contacts`; converts `idx_users_email` and `idx_profiles_user_id` into **partial** unique indexes (`WHERE deleted_at IS NULL`) so soft-deleted rows never block reuse (e.g. re-registering an email).
- **Down (best-effort):** restores plain unique indexes — fails if soft-deleted rows would violate them (hard-purge first); drops the columns.
