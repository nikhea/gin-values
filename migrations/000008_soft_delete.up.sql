-- Soft delete: rows are hidden via deleted_at instead of removed.
-- Unique indexes become partial so a deleted row never blocks reuse
-- (e.g. re-registering the same email).
ALTER TABLE users ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMPTZ;
ALTER TABLE profiles ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMPTZ;
ALTER TABLE contacts ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMPTZ;

CREATE INDEX IF NOT EXISTS idx_users_deleted_at ON users (deleted_at);
CREATE INDEX IF NOT EXISTS idx_profiles_deleted_at ON profiles (deleted_at);
CREATE INDEX IF NOT EXISTS idx_contacts_deleted_at ON contacts (deleted_at);

DROP INDEX IF EXISTS idx_users_email;
CREATE UNIQUE INDEX idx_users_email ON users (email) WHERE deleted_at IS NULL;

DROP INDEX IF EXISTS idx_profiles_user_id;
CREATE UNIQUE INDEX idx_profiles_user_id ON profiles (user_id) WHERE deleted_at IS NULL;
