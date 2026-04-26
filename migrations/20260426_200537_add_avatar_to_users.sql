-- +migrate Up
ALTER TABLE users
ADD COLUMN avatar VARCHAR(255) NULL;

-- +migrate Down
ALTER TABLE users
DROP COLUMN IF EXISTS avatar;
