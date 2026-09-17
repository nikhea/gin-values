DROP INDEX IF EXISTS idx_profiles_user_id;
ALTER TABLE IF EXISTS profiles DROP CONSTRAINT IF EXISTS fk_profiles_user;
DROP TABLE IF EXISTS profiles;
