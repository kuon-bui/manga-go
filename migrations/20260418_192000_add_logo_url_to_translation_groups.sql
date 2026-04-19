-- +migrate Up
ALTER TABLE translation_groups ADD COLUMN logo_url VARCHAR(500) NULL;

-- +migrate Down
ALTER TABLE translation_groups DROP COLUMN IF EXISTS logo_url;
