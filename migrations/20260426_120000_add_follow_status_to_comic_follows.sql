-- +migrate Up
ALTER TABLE comic_follows
ADD COLUMN follow_status VARCHAR(20) NOT NULL DEFAULT 'reading';

-- +migrate Down
ALTER TABLE comic_follows
DROP COLUMN IF EXISTS follow_status;
