-- +migrate Up
CREATE TABLE reading_histories (
    id uuid NOT NULL DEFAULT uuid_generate_v4() PRIMARY KEY,
    user_id uuid NOT NULL,
    chapter_id uuid NOT NULL,
    comic_id uuid NOT NULL,
    last_read_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ NULL,
    CONSTRAINT fk_reading_histories_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_reading_histories_chapter FOREIGN KEY (chapter_id) REFERENCES chapters(id) ON DELETE CASCADE,
    CONSTRAINT fk_reading_histories_comic FOREIGN KEY (comic_id) REFERENCES comics(id) ON DELETE CASCADE,
    CONSTRAINT uq_reading_histories_user_chapter UNIQUE(user_id, chapter_id)
);

CREATE INDEX idx_reading_histories_user_id ON reading_histories(user_id);
CREATE INDEX idx_reading_histories_comic_id ON reading_histories(comic_id);
CREATE INDEX idx_reading_histories_chapter_id ON reading_histories(chapter_id);
CREATE INDEX idx_reading_histories_last_read_at ON reading_histories(last_read_at DESC);

CREATE TRIGGER update_reading_histories_updated_at
BEFORE UPDATE ON reading_histories
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- +migrate Down
DROP TRIGGER IF EXISTS update_reading_histories_updated_at ON reading_histories;
DROP INDEX IF EXISTS idx_reading_histories_last_read_at;
DROP INDEX IF EXISTS idx_reading_histories_chapter_id;
DROP INDEX IF EXISTS idx_reading_histories_comic_id;
DROP INDEX IF EXISTS idx_reading_histories_user_id;
DROP TABLE reading_histories;
