-- +migrate Up
CREATE TABLE comic_genres (
    comic_id uuid NOT NULL REFERENCES comics(id) ON DELETE CASCADE,
    genre_id uuid NOT NULL REFERENCES genres(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (comic_id, genre_id)
);

CREATE INDEX idx_comic_genres_genre_id ON comic_genres (genre_id);

-- +migrate Down
DROP INDEX IF EXISTS idx_comic_genres_genre_id;
DROP TABLE comic_genres;
