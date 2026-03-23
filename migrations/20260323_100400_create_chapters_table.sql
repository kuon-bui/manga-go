-- +migrate Up
CREATE TABLE chapters (
    id uuid NOT NULL DEFAULT uuid_generate_v4() PRIMARY KEY,
    comic_id uuid NOT NULL REFERENCES comics(id) ON DELETE CASCADE,
    number INT NOT NULL,
    title VARCHAR(255) NOT NULL DEFAULT '',
    slug VARCHAR(255) NOT NULL DEFAULT '',
    is_published BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ NULL
);

CREATE UNIQUE INDEX idx_chapters_comic_number ON chapters (comic_id, number) WHERE deleted_at IS NULL;
CREATE UNIQUE INDEX idx_chapters_comic_slug ON chapters (comic_id, slug) WHERE deleted_at IS NULL;
CREATE INDEX idx_chapters_comic_id ON chapters (comic_id);

CREATE TRIGGER update_chapters_updated_at
BEFORE UPDATE ON chapters
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- +migrate Down
DROP TRIGGER IF EXISTS update_chapters_updated_at ON chapters;
DROP INDEX IF EXISTS idx_chapters_comic_id;
DROP INDEX IF EXISTS idx_chapters_comic_slug;
DROP INDEX IF EXISTS idx_chapters_comic_number;
DROP TABLE chapters;
