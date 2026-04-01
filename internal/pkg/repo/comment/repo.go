package commentrepo

import (
	"manga-go/internal/pkg/model"
	"manga-go/internal/pkg/repo/base"

	"gorm.io/gorm"
)

type CommentRepo struct {
	*base.BaseRepository[model.Comment]
}

func NewCommentRepo(db *gorm.DB) *CommentRepo {
	return &CommentRepo{
		BaseRepository: &base.BaseRepository[model.Comment]{
			DB: db,
		},
	}
}
