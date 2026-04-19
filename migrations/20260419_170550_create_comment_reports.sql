-- +migrate Up
CREATE TABLE comment_reports (
    id uuid NOT NULL DEFAULT uuid_generate_v4(),
    comment_id uuid NOT NULL REFERENCES comments(id) ON DELETE CASCADE,
    user_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    reason VARCHAR(50) NOT NULL,
    details TEXT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ NULL,
    PRIMARY KEY (id)
);

CREATE INDEX idx_comment_reports_comment_id ON comment_reports(comment_id);
CREATE INDEX idx_comment_reports_user_id ON comment_reports(user_id);
CREATE INDEX idx_comment_reports_deleted_at ON comment_reports(deleted_at);

CREATE TRIGGER update_comment_reports_updated_at
BEFORE UPDATE ON comment_reports
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- +migrate Down
DROP TRIGGER IF EXISTS update_comment_reports_updated_at ON comment_reports;
DROP INDEX IF EXISTS idx_comment_reports_deleted_at;
DROP INDEX IF EXISTS idx_comment_reports_user_id;
DROP INDEX IF EXISTS idx_comment_reports_comment_id;
DROP TABLE comment_reports;
