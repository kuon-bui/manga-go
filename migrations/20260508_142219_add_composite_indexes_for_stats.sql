-- +migrate Up notransaction
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_comic_follows_comic_deleted_at
ON comic_follows(comic_id, deleted_at);

CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_ratings_comic_deleted_at
ON ratings(comic_id, deleted_at);

-- chapters đã có partial index (comic_id, number) WHERE deleted_at IS NULL,
-- nhưng thiếu index cho (comic_id, deleted_at) để optimize subquery
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_chapters_comic_deleted_at
ON chapters(comic_id, deleted_at);

-- +migrate Down
DROP INDEX IF EXISTS idx_chapters_comic_deleted_at;
DROP INDEX IF EXISTS idx_ratings_comic_deleted_at;
DROP INDEX IF EXISTS idx_comic_follows_comic_deleted_at;
