ALTER TABLE users
    DROP COLUMN IF EXISTS otp_attempts,
    DROP COLUMN IF EXISTS otp_expires_at,
    DROP COLUMN IF EXISTS otp_hash;
