-- +migrate Up
ALTER TABLE users
ADD COLUMN user_config BYTEA NOT NULL DEFAULT decode('0a', 'hex');

-- +migrate Down
ALTER TABLE users
DROP COLUMN IF EXISTS user_config;
