-- +migrate Up
CREATE TABLE comic_stats (
    comic_id uuid NOT NULL PRIMARY KEY REFERENCES comics(id) ON DELETE CASCADE,
    follow_count INT NOT NULL DEFAULT 0,
    rating_count INT NOT NULL DEFAULT 0,
    chapter_count INT NOT NULL DEFAULT 0,
    avg_rating FLOAT NULL,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- +migrate Down
DROP TABLE comic_stats;
