-- +migrate Up
CREATE TABLE comment_reactions (
	id uuid NOT NULL DEFAULT uuid_generate_v4() PRIMARY KEY,
	user_id uuid NOT NULL,
	comment_id uuid NOT NULL,
	type varchar(255) NOT NULL,
	created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
	deleted_at TIMESTAMPTZ NULL,
	CONSTRAINT fk_comment_reactions_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
	CONSTRAINT fk_comment_reactions_comment FOREIGN KEY (comment_id) REFERENCES comments(id) ON DELETE CASCADE
);

CREATE UNIQUE INDEX idx_comment_reactions_user_comment ON comment_reactions(user_id, comment_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_comment_reactions_user_id ON comment_reactions(user_id);
CREATE INDEX idx_comment_reactions_comment_id ON comment_reactions(comment_id);
CREATE INDEX idx_comment_reactions_comment_type ON comment_reactions(comment_id, type) WHERE deleted_at IS NULL;

CREATE TRIGGER update_comment_reactions_updated_at
BEFORE UPDATE ON comment_reactions
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- +migrate Down
DROP TRIGGER IF EXISTS update_comment_reactions_updated_at ON comment_reactions;
DROP INDEX IF EXISTS idx_comment_reactions_comment_type;
DROP INDEX IF EXISTS idx_comment_reactions_comment_id;
DROP INDEX IF EXISTS idx_comment_reactions_user_id;
DROP INDEX IF EXISTS idx_comment_reactions_user_comment;
DROP TABLE IF EXISTS comment_reactions;
