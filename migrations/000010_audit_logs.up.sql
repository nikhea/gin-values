-- Append-only audit trail. Rows are never updated or deleted by the app.
CREATE TABLE IF NOT EXISTS audit_logs (
    id VARCHAR(36) PRIMARY KEY,
    actor_id VARCHAR(36),
    action VARCHAR(64) NOT NULL,
    resource_type VARCHAR(32) NOT NULL,
    resource_id VARCHAR(36),
    org_id VARCHAR(36),
    ip TEXT,
    request_id VARCHAR(36),
    metadata JSONB NOT NULL DEFAULT '{}',
    created_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_audit_logs_actor_id ON audit_logs (actor_id);
CREATE INDEX IF NOT EXISTS idx_audit_logs_org_id ON audit_logs (org_id);
CREATE INDEX IF NOT EXISTS idx_audit_logs_action ON audit_logs (action);
CREATE INDEX IF NOT EXISTS idx_audit_logs_created_at ON audit_logs (created_at);
