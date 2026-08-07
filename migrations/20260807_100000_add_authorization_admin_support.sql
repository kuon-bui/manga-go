-- +migrate Up
ALTER TABLE roles ADD COLUMN description TEXT NULL;

CREATE TABLE authorization_cache_revisions (
    scope VARCHAR(128) PRIMARY KEY,
    version BIGINT NOT NULL DEFAULT 1,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO authorization_cache_revisions (scope, version)
VALUES ('global', 1);

CREATE TABLE authorization_audit_logs (
    id uuid NOT NULL DEFAULT uuid_generate_v4() PRIMARY KEY,
    actor_user_id uuid NULL,
    actor_name_snapshot VARCHAR(255) NOT NULL,
    actor_email_snapshot VARCHAR(255) NOT NULL,
    action VARCHAR(64) NOT NULL,
    target_type VARCHAR(32) NOT NULL,
    target_id uuid NOT NULL,
    target_name_snapshot VARCHAR(255) NOT NULL,
    before JSONB NOT NULL DEFAULT '{}'::jsonb,
    after JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_authorization_audit_logs_created_at
    ON authorization_audit_logs(created_at DESC);
CREATE INDEX idx_authorization_audit_logs_actor
    ON authorization_audit_logs(actor_user_id, created_at DESC);
CREATE INDEX idx_authorization_audit_logs_target
    ON authorization_audit_logs(target_type, target_id, created_at DESC);

-- +migrate Down
DROP TABLE IF EXISTS authorization_audit_logs;
DROP TABLE IF EXISTS authorization_cache_revisions;
ALTER TABLE roles DROP COLUMN IF EXISTS description;
