-- +migrate Up
ALTER TABLE comics
    ADD COLUMN translation_group_id uuid NULL REFERENCES translation_groups(id);

-- +migrate Down
ALTER TABLE comics DROP COLUMN IF EXISTS translation_group_id;
