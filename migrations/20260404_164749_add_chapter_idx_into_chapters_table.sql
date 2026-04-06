-- +migrate Up
ALTER TABLE chapters ADD COLUMN chapter_idx INT NOT NULL DEFAULT 0;

-- +migrate Down
ALTER TABLE chapters DROP COLUMN chapter_idx;
