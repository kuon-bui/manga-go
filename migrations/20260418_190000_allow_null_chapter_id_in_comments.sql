-- +migrate Up
-- Allow comic-level comments (not tied to a specific chapter)
ALTER TABLE comments ALTER COLUMN chapter_id DROP NOT NULL;

-- Drop the old unique constraint that assumed chapter_id non-null
ALTER TABLE comments DROP CONSTRAINT IF EXISTS uq_comments_user_chapter;

-- Partial index: unique constraint only when chapter_id is not null (per-chapter scope)
-- CREATE UNIQUE INDEX IF NOT EXISTS uq_comments_user_chapter_scoped
--     ON comments (user_id, chapter_id)
--     WHERE chapter_id IS NOT NULL AND deleted_at IS NULL;

-- Index to efficiently query comic-level comments (chapter_id IS NULL)
-- CREATE INDEX IF NOT EXISTS idx_comments_comic_level
--     ON comments (comic_id)
--     WHERE chapter_id IS NULL AND parent_id IS NULL AND deleted_at IS NULL;

-- +migrate Down
DROP INDEX IF EXISTS idx_comments_comic_level;
DROP INDEX IF EXISTS uq_comments_user_chapter_scoped;
ALTER TABLE comments ALTER COLUMN chapter_id SET NOT NULL;
