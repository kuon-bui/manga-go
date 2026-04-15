-- +migrate Up
CREATE TABLE notifications (
    id uuid NOT NULL DEFAULT uuid_generate_v4() PRIMARY KEY,
    type VARCHAR(100) NOT NULL,
    category VARCHAR(50) NOT NULL,
    actor_id uuid NULL,
    entity_type VARCHAR(50) NULL,
    entity_id uuid NULL,
    dedupe_key VARCHAR(255) NULL,
    title VARCHAR(255) NOT NULL,
    body TEXT NOT NULL,
    payload JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ NULL,
    CONSTRAINT fk_notifications_actor FOREIGN KEY (actor_id) REFERENCES users(id)
);

CREATE UNIQUE INDEX ux_notifications_dedupe_key
ON notifications (dedupe_key)
WHERE dedupe_key IS NOT NULL AND deleted_at IS NULL;

CREATE INDEX idx_notifications_type_created_at
ON notifications (type, created_at DESC)
WHERE deleted_at IS NULL;

CREATE TRIGGER update_notifications_updated_at
BEFORE UPDATE ON notifications
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- +migrate Down
DROP TRIGGER IF EXISTS update_notifications_updated_at ON notifications;
DROP INDEX IF EXISTS idx_notifications_type_created_at;
DROP INDEX IF EXISTS ux_notifications_dedupe_key;
DROP TABLE IF EXISTS notifications;
