-- +migrate Up
ALTER TABLE comments DROP CONSTRAINT IF EXISTS uq_comments_user_chapter;

ALTER TABLE comments
ADD COLUMN parent_id uuid NULL;

ALTER TABLE comments
ADD CONSTRAINT fk_comments_parent FOREIGN KEY (parent_id) REFERENCES comments(id) ON DELETE CASCADE;

ALTER TABLE comments
ADD CONSTRAINT chk_comments_parent_not_self CHECK (parent_id IS NULL OR parent_id <> id);

CREATE INDEX idx_comments_parent_id ON comments(parent_id);
CREATE INDEX idx_comments_parent_id_created_at ON comments(parent_id, created_at);

-- +migrate Down
DROP INDEX IF EXISTS idx_comments_parent_id_created_at;
DROP INDEX IF EXISTS idx_comments_parent_id;

ALTER TABLE comments DROP CONSTRAINT IF EXISTS chk_comments_parent_not_self;
ALTER TABLE comments DROP CONSTRAINT IF EXISTS fk_comments_parent;
ALTER TABLE comments DROP COLUMN IF EXISTS parent_id;

ALTER TABLE comments
ADD CONSTRAINT uq_comments_user_chapter UNIQUE(user_id, chapter_id);
