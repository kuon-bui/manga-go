-- +migrate Up
CREATE TABLE permissions (
    id uuid NOT NULL DEFAULT uuid_generate_v4() PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ NULL
);

CREATE TRIGGER update_permissions_updated_at
BEFORE UPDATE ON permissions
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- +migrate Down
DROP TRIGGER IF EXISTS update_permissions_updated_at ON permissions;
DROP TABLE permissions;
