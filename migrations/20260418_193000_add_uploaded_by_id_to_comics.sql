-- +migrate Up
ALTER TABLE comics ADD COLUMN uploaded_by_id UUID NULL REFERENCES users(id) ON DELETE SET NULL;

-- +migrate Down
ALTER TABLE comics DROP COLUMN IF EXISTS uploaded_by_id;
