-- +migrate Up
CREATE TABLE comics (
    id uuid NOT NULL DEFAULT uuid_generate_v4() PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    slug VARCHAR(255) NOT NULL,
    alternative_titles JSONB NULL,
    description TEXT NULL,
    thumbnail VARCHAR(500) NULL,
    banner VARCHAR(500) NULL,
    type varchar(50) NOT NULL DEFAULT 'manga',
    status varchar(50) NOT NULL DEFAULT 'ongoing',
    is_active BOOLEAN NOT NULL DEFAULT true,
    is_hot BOOLEAN NOT NULL DEFAULT false,
    is_featured BOOLEAN NOT NULL DEFAULT false,
    author VARCHAR(255) NULL,
    artist VARCHAR(255) NULL,
    published_year SMALLINT NULL,
    last_chapter_at TIMESTAMPTZ NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ NULL
);

CREATE UNIQUE INDEX idx_comics_slug ON comics (slug) WHERE deleted_at IS NULL;

CREATE TRIGGER update_comics_updated_at
BEFORE UPDATE ON comics
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- +migrate Down
DROP TRIGGER IF EXISTS update_comics_updated_at ON comics;
DROP INDEX IF EXISTS idx_comics_slug;
DROP TABLE comics;
DROP TYPE IF EXISTS comic_status;
DROP TYPE IF EXISTS comic_type;
