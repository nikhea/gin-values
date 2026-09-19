DROP INDEX IF EXISTS idx_casbin_rule_ptype_v0;
DROP TABLE IF EXISTS casbin_rule;

DROP INDEX IF EXISTS idx_contacts_org_id;
ALTER TABLE contacts DROP COLUMN IF EXISTS org_id;

DROP INDEX IF EXISTS idx_memberships_org_id;
DROP TABLE IF EXISTS memberships;
DROP TABLE IF EXISTS organizations;
