-- +migrate Up
CREATE TABLE page_reactions (
	id uuid NOT NULL DEFAULT uuid_generate_v4() PRIMARY KEY,
	user_id uuid NOT NULL,
	page_id uuid NOT NULL,
	type varchar(255) NOT NULL,
	created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
	deleted_at TIMESTAMPTZ NULL,
	CONSTRAINT fk_page_reactions_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
	CONSTRAINT fk_page_reactions_page FOREIGN KEY (page_id) REFERENCES pages(id) ON DELETE CASCADE
);

CREATE UNIQUE INDEX idx_page_reactions_user_page ON page_reactions(user_id, page_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_page_reactions_user_id ON page_reactions(user_id);
CREATE INDEX idx_page_reactions_page_id ON page_reactions(page_id);
CREATE INDEX idx_page_reactions_page_type ON page_reactions(page_id, type) WHERE deleted_at IS NULL;

CREATE TRIGGER update_page_reactions_updated_at
BEFORE UPDATE ON page_reactions
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- +migrate Down
DROP TRIGGER IF EXISTS update_page_reactions_updated_at ON page_reactions;
DROP INDEX IF EXISTS idx_page_reactions_page_type;
DROP INDEX IF EXISTS idx_page_reactions_page_id;
DROP INDEX IF EXISTS idx_page_reactions_user_id;
DROP INDEX IF EXISTS idx_page_reactions_user_page;
DROP TABLE IF EXISTS page_reactions;
