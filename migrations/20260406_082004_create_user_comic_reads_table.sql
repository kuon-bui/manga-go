-- +migrate Up
CREATE TABLE user_comic_reads (
    id uuid NOT NULL DEFAULT uuid_generate_v4() PRIMARY KEY,
    user_id uuid NOT NULL,
    comic_id uuid NOT NULL,
    read_data BYTEA NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ NULL,
    CONSTRAINT fk_user_comic_reads_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_user_comic_reads_comic FOREIGN KEY (comic_id) REFERENCES comics(id) ON DELETE CASCADE
);

CREATE INDEX idx_user_comic_reads_user_comic_id ON user_comic_reads(user_id, comic_id);

CREATE TRIGGER update_user_comic_reads_updated_at
BEFORE UPDATE ON user_comic_reads
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- +migrate Down
DROP TRIGGER IF EXISTS update_user_comic_reads_updated_at ON user_comic_reads;
DROP INDEX IF EXISTS idx_user_comic_reads_user_comic_id;
DROP TABLE user_comic_reads;
