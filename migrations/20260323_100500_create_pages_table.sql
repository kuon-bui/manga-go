-- +migrate Up
CREATE TABLE pages (
    id uuid NOT NULL DEFAULT uuid_generate_v4() PRIMARY KEY,
    chapter_id uuid NOT NULL REFERENCES chapters(id) ON DELETE CASCADE,
    page_number INT NOT NULL,
    image_url VARCHAR(500) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ NULL
);

CREATE UNIQUE INDEX idx_pages_chapter_page_number ON pages (chapter_id, page_number) WHERE deleted_at IS NULL;
CREATE INDEX idx_pages_chapter_id ON pages (chapter_id);

CREATE TRIGGER update_pages_updated_at
BEFORE UPDATE ON pages
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- +migrate Down
DROP TRIGGER IF EXISTS update_pages_updated_at ON pages;
DROP INDEX IF EXISTS idx_pages_chapter_id;
DROP INDEX IF EXISTS idx_pages_chapter_page_number;
DROP TABLE pages;
