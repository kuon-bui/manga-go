-- +migrate Up
ALTER TABLE comics RENAME COLUMN is_active TO is_published;

-- +migrate Down
ALTER TABLE comics RENAME COLUMN is_published TO is_active;
