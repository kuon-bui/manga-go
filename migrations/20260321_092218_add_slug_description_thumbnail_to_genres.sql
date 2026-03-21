-- +migrate Up
ALTER TABLE genres
    ADD COLUMN slug VARCHAR(255) NOT NULL DEFAULT '',
    ADD COLUMN description TEXT NOT NULL DEFAULT '',
    ADD COLUMN thumbnail VARCHAR(255) NOT NULL DEFAULT '';

-- +migrate Down
ALTER TABLE genres
    DROP COLUMN IF EXISTS slug,
    DROP COLUMN IF EXISTS description,
    DROP COLUMN IF EXISTS thumbnail;
