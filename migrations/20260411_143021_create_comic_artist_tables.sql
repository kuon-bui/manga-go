-- +migrate Up

CREATE TABLE comic_artists (
    comic_id uuid NOT NULL REFERENCES comics(id) ON DELETE CASCADE,
    artist_id uuid NOT NULL REFERENCES authors(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (comic_id, artist_id)
);


-- +migrate Down
DROP TABLE IF EXISTS comic_artists;
