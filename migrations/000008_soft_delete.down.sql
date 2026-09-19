-- Best-effort rollback: fails if soft-deleted rows would violate the
-- restored plain unique indexes (delete or hard-purge them first).
DROP INDEX IF EXISTS idx_profiles_user_id;
CREATE UNIQUE INDEX idx_profiles_user_id ON profiles (user_id);

DROP INDEX IF EXISTS idx_users_email;
CREATE UNIQUE INDEX idx_users_email ON users (email);

DROP INDEX IF EXISTS idx_contacts_deleted_at;
DROP INDEX IF EXISTS idx_profiles_deleted_at;
DROP INDEX IF EXISTS idx_users_deleted_at;

ALTER TABLE contacts DROP COLUMN IF EXISTS deleted_at;
ALTER TABLE profiles DROP COLUMN IF EXISTS deleted_at;
ALTER TABLE users DROP COLUMN IF EXISTS deleted_at;
