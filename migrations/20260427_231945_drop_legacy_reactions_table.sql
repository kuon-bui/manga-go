-- +migrate Up
DROP TRIGGER IF EXISTS update_reactions_updated_at ON reactions;
DROP INDEX IF EXISTS idx_reactions_comment_id;
DROP INDEX IF EXISTS idx_reactions_user_id;
DROP TABLE IF EXISTS reactions;

-- +migrate Down
CREATE TABLE reactions (
	id uuid NOT NULL DEFAULT uuid_generate_v4() PRIMARY KEY,
	user_id uuid NOT NULL,
	comment_id uuid NOT NULL,
	type varchar(255) NOT NULL,
	last_read_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
	created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
	deleted_at TIMESTAMPTZ NULL,
	CONSTRAINT fk_reactions_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
	CONSTRAINT fk_reactions_comment FOREIGN KEY (comment_id) REFERENCES comments(id) ON DELETE CASCADE
);

CREATE INDEX idx_reactions_user_id ON reactions(user_id);
CREATE INDEX idx_reactions_comment_id ON reactions(comment_id);

CREATE TRIGGER update_reactions_updated_at
BEFORE UPDATE ON reactions
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
