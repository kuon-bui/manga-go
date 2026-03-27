-- +migrate Up
CREATE TABLE comic_tags (
    comic_id uuid NOT NULL REFERENCES comics(id) ON DELETE CASCADE,
    tag_id uuid NOT NULL REFERENCES tags(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (comic_id, tag_id)
);

CREATE INDEX idx_comic_tags_tag_id ON comic_tags (tag_id);

-- +migrate Down
DROP INDEX IF EXISTS idx_comic_tags_tag_id;
DROP TABLE comic_tags;
