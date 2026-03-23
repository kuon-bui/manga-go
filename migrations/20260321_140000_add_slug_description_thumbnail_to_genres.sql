-- +migrate Up
ALTER TABLE genres ADD COLUMN slug VARCHAR(255) NOT NULL DEFAULT '';
ALTER TABLE genres ADD COLUMN description TEXT NOT NULL DEFAULT '';
ALTER TABLE genres ADD COLUMN thumbnail VARCHAR(500) NOT NULL DEFAULT '';

CREATE UNIQUE INDEX idx_genres_slug ON genres (slug) WHERE deleted_at IS NULL;

-- +migrate Down
DROP INDEX IF EXISTS idx_genres_slug;
ALTER TABLE genres DROP COLUMN IF EXISTS thumbnail;
ALTER TABLE genres DROP COLUMN IF EXISTS description;
ALTER TABLE genres DROP COLUMN IF EXISTS slug;
