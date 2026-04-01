-- +migrate Up
ALTER TABLE comics
ADD COLUMN age_rating VARCHAR(10) NOT NULL DEFAULT 'all';

-- +migrate Down
ALTER TABLE comics
DROP COLUMN IF EXISTS age_rating;
