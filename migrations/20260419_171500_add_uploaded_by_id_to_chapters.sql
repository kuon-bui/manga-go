-- +migrate Up
ALTER TABLE chapters ADD COLUMN uploaded_by_id uuid;

ALTER TABLE chapters ADD CONSTRAINT fk_chapters_uploaded_by_id
FOREIGN KEY (uploaded_by_id) REFERENCES users(id) ON DELETE SET NULL;

CREATE INDEX idx_chapters_uploaded_by_id ON chapters(uploaded_by_id);

-- +migrate Down
DROP INDEX IF EXISTS idx_chapters_uploaded_by_id;
ALTER TABLE chapters DROP CONSTRAINT IF EXISTS fk_chapters_uploaded_by_id;
ALTER TABLE chapters DROP COLUMN IF EXISTS uploaded_by_id;
