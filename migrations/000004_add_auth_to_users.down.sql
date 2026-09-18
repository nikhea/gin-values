DROP INDEX IF EXISTS idx_users_reset_token;
DROP INDEX IF EXISTS idx_users_verification_token;

ALTER TABLE users
    DROP COLUMN IF EXISTS reset_expires_at,
    DROP COLUMN IF EXISTS reset_token,
    DROP COLUMN IF EXISTS verification_token,
    DROP COLUMN IF EXISTS email_verified,
    DROP COLUMN IF EXISTS password_hash;
