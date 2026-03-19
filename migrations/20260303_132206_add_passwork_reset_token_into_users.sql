-- +migrate Up
ALTER TABLE users ADD COLUMN reset_password_token VARCHAR(40) NULL,
ADD COLUMN reset_password_expiry_at TIMESTAMPTZ NULL;

-- +migrate Down
ALTER TABLE users DROP COLUMN IF EXISTS reset_password_token,
DROP COLUMN IF EXISTS reset_password_expiry_at;
