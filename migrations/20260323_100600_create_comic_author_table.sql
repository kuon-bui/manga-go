-- +migrate Up
ALTER TABLE comics DROP COLUMN IF EXISTS author;

CREATE TABLE comic_authors (
  comic_id uuid NOT NULL REFERENCES comics(id) ON DELETE CASCADE,
  author_id uuid NOT NULL REFERENCES authors(id) ON DELETE CASCADE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (comic_id, author_id)
);

CREATE INDEX idx_comic_authors_author_id ON comic_authors (author_id);

-- +migrate Down
DROP INDEX IF EXISTS idx_comic_authors_author_id;
DROP TABLE comic_authors;
ALTER TABLE comics ADD COLUMN author VARCHAR(255) NULL;
