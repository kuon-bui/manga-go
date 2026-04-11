-- +migrate Up
ALTER TABLE comics
DROP COLUMN IF EXISTS artist_id;
-- +migrate Down
ALTER TABLE comics
ADD COLUMN artist_id uuid REFERENCES authors(id) ON DELETE SET NULL;
