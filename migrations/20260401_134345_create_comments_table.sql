-- +migrate Up
CREATE TABLE comments (
    id uuid NOT NULL DEFAULT uuid_generate_v4() PRIMARY KEY,
    user_id uuid NOT NULL,
    chapter_id uuid NOT NULL,
    comic_id uuid NOT NULL,
    page_index integer NULL,
    content text NOT NULL,
    last_read_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ NULL,
    CONSTRAINT fk_comments_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_comments_chapter FOREIGN KEY (chapter_id) REFERENCES chapters(id) ON DELETE CASCADE,
    CONSTRAINT fk_comments_comic FOREIGN KEY (comic_id) REFERENCES comics(id) ON DELETE CASCADE,
    CONSTRAINT uq_comments_user_chapter UNIQUE(user_id, chapter_id)
);
CREATE INDEX idx_comments_user_id ON comments(user_id);
CREATE INDEX idx_comments_comic_id ON comments(comic_id);
CREATE INDEX idx_comments_chapter_id ON comments(chapter_id);
CREATE INDEX idx_comments_page_index ON comments(chapter_id, page_index);


CREATE TRIGGER update_comments_updated_at
BEFORE UPDATE ON comments
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- +migrate Down
DROP TRIGGER IF EXISTS update_comments_updated_at ON comments;
DROP INDEX IF EXISTS idx_comments_page_index;
DROP INDEX IF EXISTS idx_comments_chapter_id;
DROP INDEX IF EXISTS idx_comments_comic_id;
DROP INDEX IF EXISTS idx_comments_user_id; 
DROP TABLE comments;
