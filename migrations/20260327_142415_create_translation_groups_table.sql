-- +migrate Up
CREATE TABLE translation_groups (
    id uuid NOT NULL DEFAULT uuid_generate_v4() PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    owner_id uuid NOT NULL REFERENCES users(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ NULL
);

CREATE TRIGGER update_translation_groups_updated_at
BEFORE UPDATE ON translation_groups
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

ALTER TABLE users
    ADD COLUMN translation_group_id uuid NULL REFERENCES translation_groups(id);

-- +migrate Down
ALTER TABLE users DROP COLUMN IF EXISTS translation_group_id;

DROP TRIGGER IF EXISTS update_translation_groups_updated_at ON translation_groups;
DROP TABLE translation_groups;
