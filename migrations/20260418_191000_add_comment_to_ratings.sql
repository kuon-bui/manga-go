-- +migrate Up
ALTER TABLE ratings ADD COLUMN comment TEXT NULL;

-- +migrate Down
ALTER TABLE ratings DROP COLUMN IF EXISTS comment;
