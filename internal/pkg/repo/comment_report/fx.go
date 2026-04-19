package commentreportrepo

import (
	"go.uber.org/fx"
	"gorm.io/gorm"
)

var Module = fx.Module(
	"comment-report-repo",
	fx.Provide(NewCommentReportRepo),
)

func Provide(db *gorm.DB) *CommentReportRepo {
	return NewCommentReportRepo(db)
}
