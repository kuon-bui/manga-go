-- +migrate Up
ALTER TABLE comics
    ADD COLUMN translation_group_id uuid NULL REFERENCES translation_groups(id) ON DELETE SET NULL;

CREATE INDEX idx_comics_translation_group_id ON comics (translation_group_id) WHERE deleted_at IS NULL;

-- +migrate Down
DROP INDEX IF EXISTS idx_comics_translation_group_id;

ALTER TABLE comics
    DROP COLUMN IF EXISTS translation_group_id;