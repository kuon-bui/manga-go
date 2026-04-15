-- +migrate Up
CREATE TABLE comic_follows (
    id uuid NOT NULL DEFAULT uuid_generate_v4() PRIMARY KEY,
    user_id uuid NOT NULL,
    comic_id uuid NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ NULL,
    CONSTRAINT fk_comic_follows_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_comic_follows_comic FOREIGN KEY (comic_id) REFERENCES comics(id) ON DELETE CASCADE,
    CONSTRAINT uq_comic_follows_user_comic UNIQUE (user_id, comic_id)
);

CREATE INDEX idx_comic_follows_user_id ON comic_follows(user_id);
CREATE INDEX idx_comic_follows_comic_id ON comic_follows(comic_id);

CREATE TRIGGER update_comic_follows_updated_at
BEFORE UPDATE ON comic_follows
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- +migrate Down
DROP TRIGGER IF EXISTS update_comic_follows_updated_at ON comic_follows;
DROP INDEX IF EXISTS idx_comic_follows_comic_id;
DROP INDEX IF EXISTS idx_comic_follows_user_id;
DROP TABLE comic_follows;
