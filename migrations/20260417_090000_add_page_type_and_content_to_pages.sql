-- +migrate Up
ALTER TABLE pages
    ADD COLUMN IF NOT EXISTS page_type VARCHAR(20);

UPDATE pages
SET page_type = 'image'
WHERE page_type IS NULL;

ALTER TABLE pages
    ALTER COLUMN page_type SET DEFAULT 'image',
    ALTER COLUMN page_type SET NOT NULL;

ALTER TABLE pages
    ADD COLUMN IF NOT EXISTS content TEXT;

ALTER TABLE pages
    ALTER COLUMN image_url DROP NOT NULL;

ALTER TABLE pages
    ADD CONSTRAINT chk_pages_page_type CHECK (page_type IN ('image', 'text'));

-- +migrate Down
ALTER TABLE pages
    DROP CONSTRAINT IF EXISTS chk_pages_page_type;

ALTER TABLE pages
    DROP COLUMN IF EXISTS content,
    DROP COLUMN IF EXISTS page_type;

UPDATE pages
SET image_url = ''
WHERE image_url IS NULL;

ALTER TABLE pages
    ALTER COLUMN image_url SET NOT NULL;
