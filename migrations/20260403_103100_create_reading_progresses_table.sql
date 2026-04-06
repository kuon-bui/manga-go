-- +migrate Up
CREATE TABLE reading_progresses (
    id uuid NOT NULL DEFAULT uuid_generate_v4() PRIMARY KEY,
    user_id uuid NOT NULL,
    comic_id uuid NOT NULL,
    chapter_id uuid NOT NULL,
    scroll_percent int NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ NULL,
    CONSTRAINT fk_reading_progresses_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_reading_progresses_comic FOREIGN KEY (comic_id) REFERENCES comics(id) ON DELETE CASCADE,
    CONSTRAINT fk_reading_progresses_chapter FOREIGN KEY (chapter_id) REFERENCES chapters(id) ON DELETE CASCADE
);

CREATE INDEX idx_reading_progresses_user_id ON reading_progresses(user_id);
CREATE INDEX idx_reading_progresses_comic_id ON reading_progresses(comic_id);
CREATE INDEX idx_reading_progresses_chapter_id ON reading_progresses(chapter_id);

CREATE TRIGGER update_reading_progresses_updated_at
BEFORE UPDATE ON reading_progresses
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- +migrate Down
DROP TRIGGER IF EXISTS update_reading_progresses_updated_at ON reading_progresses;
DROP INDEX IF EXISTS idx_reading_progresses_chapter_id;
DROP INDEX IF EXISTS idx_reading_progresses_comic_id;
DROP INDEX IF EXISTS idx_reading_progresses_user_id;
DROP TABLE reading_progresses;
