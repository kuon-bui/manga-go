-- +migrate Up
ALTER TABLE users
ADD COLUMN user_config BYTEA NOT NULL DEFAULT decode('0a', 'hex');

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

CREATE TABLE user_notifications (
    id uuid NOT NULL DEFAULT uuid_generate_v4() PRIMARY KEY,
    notification_id uuid NOT NULL,
    user_id uuid NOT NULL,
    channel_state BIGINT NOT NULL DEFAULT 0,
    is_seen BOOLEAN NOT NULL DEFAULT FALSE,
    seen_at TIMESTAMPTZ NULL,
    is_read BOOLEAN NOT NULL DEFAULT FALSE,
    read_at TIMESTAMPTZ NULL,
    emailed_at TIMESTAMPTZ NULL,
    pushed_at TIMESTAMPTZ NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ NULL,
    CONSTRAINT fk_user_notifications_notification FOREIGN KEY (notification_id) REFERENCES notifications(id),
    CONSTRAINT fk_user_notifications_user FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE UNIQUE INDEX ux_user_notifications_notification_user
ON user_notifications (notification_id, user_id)
WHERE deleted_at IS NULL;

CREATE INDEX idx_user_notifications_user_created_at
ON user_notifications (user_id, created_at DESC)
WHERE deleted_at IS NULL;

CREATE INDEX idx_user_notifications_user_unread
ON user_notifications (user_id, is_read, is_seen, created_at DESC)
WHERE deleted_at IS NULL;

CREATE TRIGGER update_user_notifications_updated_at
BEFORE UPDATE ON user_notifications
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- +migrate Down
DROP TRIGGER IF EXISTS update_user_notifications_updated_at ON user_notifications;
DROP INDEX IF EXISTS idx_user_notifications_user_unread;
DROP INDEX IF EXISTS idx_user_notifications_user_created_at;
DROP INDEX IF EXISTS ux_user_notifications_notification_user;
DROP TABLE IF EXISTS user_notifications;

DROP TRIGGER IF EXISTS update_notifications_updated_at ON notifications;
DROP INDEX IF EXISTS idx_notifications_type_created_at;
DROP INDEX IF EXISTS ux_notifications_dedupe_key;
DROP TABLE IF EXISTS notifications;

ALTER TABLE users
DROP COLUMN IF EXISTS user_config;
