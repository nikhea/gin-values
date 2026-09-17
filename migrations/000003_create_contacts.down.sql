DROP INDEX IF EXISTS idx_contacts_type;
DROP INDEX IF EXISTS idx_contacts_user_id;
ALTER TABLE IF EXISTS contacts DROP CONSTRAINT IF EXISTS fk_contacts_user;
DROP TABLE IF EXISTS contacts;
