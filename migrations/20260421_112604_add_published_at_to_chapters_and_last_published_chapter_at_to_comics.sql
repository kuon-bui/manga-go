-- +migrate Up
ALTER TABLE chapters
ADD COLUMN published_at TIMESTAMPTZ NULL;

-- +migrate Down
ALTER TABLE chapters
DROP COLUMN IF EXISTS published_at;
