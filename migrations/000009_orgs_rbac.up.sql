CREATE TABLE IF NOT EXISTS organizations (
    id VARCHAR(36) PRIMARY KEY,
    name TEXT NOT NULL,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS memberships (
    user_id VARCHAR(36) NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    org_id VARCHAR(36) NOT NULL REFERENCES organizations (id) ON DELETE CASCADE,
    role VARCHAR(16) NOT NULL,
    created_at TIMESTAMPTZ,
    PRIMARY KEY (user_id, org_id)
);

CREATE INDEX IF NOT EXISTS idx_memberships_org_id ON memberships (org_id);

-- Nullable: NULL means a personal contact owned by user_id.
ALTER TABLE contacts ADD COLUMN IF NOT EXISTS org_id VARCHAR(36) REFERENCES organizations (id) ON DELETE SET NULL;
CREATE INDEX IF NOT EXISTS idx_contacts_org_id ON contacts (org_id);

-- Casbin policy storage (managed by gorm-adapter, created here for
-- migrate-driven environments so boot never depends on AutoMigrate).
CREATE TABLE IF NOT EXISTS casbin_rule (
    id SERIAL PRIMARY KEY,
    ptype VARCHAR(100),
    v0 VARCHAR(100),
    v1 VARCHAR(100),
    v2 VARCHAR(100),
    v3 VARCHAR(100),
    v4 VARCHAR(100),
    v5 VARCHAR(100)
);
CREATE INDEX IF NOT EXISTS idx_casbin_rule_ptype_v0 ON casbin_rule (ptype, v0);
