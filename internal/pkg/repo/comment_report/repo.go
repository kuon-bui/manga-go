package commentreportrepo

import (
	"manga-go/internal/pkg/model"
	"manga-go/internal/pkg/repo/base"

	"gorm.io/gorm"
)

type CommentReportRepo struct {
	*base.BaseRepository[model.CommentReport]
}

func NewCommentReportRepo(db *gorm.DB) *CommentReportRepo {
	return &CommentReportRepo{
		BaseRepository: &base.BaseRepository[model.CommentReport]{
			DB: db,
		},
	}
}
