package commentrepo

import (
	"context"
	"manga-go/internal/pkg/common"
	"manga-go/internal/pkg/model"

	"gorm.io/gorm"
)

func (r *CommentRepo) FindRecentTopLevelPaginated(ctx context.Context, paging *common.Paging) ([]*model.Comment, int64, error) {
	var comments []*model.Comment
	var total int64

	db := r.DB.WithContext(ctx).
		Model(&model.Comment{}).
		Where("parent_id IS NULL")

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	query := r.DB.WithContext(ctx).
		Model(&model.Comment{}).
		Where("parent_id IS NULL").
		Order("comments.created_at DESC").
		Preload("User", func(tx *gorm.DB) *gorm.DB {
			return tx.Select("id", "name", "avatar")
		}).
		Preload("Comic", func(tx *gorm.DB) *gorm.DB {
			return tx.Select("id", "title", "slug", "status")
		}).
		Preload("Chapter", func(tx *gorm.DB) *gorm.DB {
			return tx.Select("id", "title", "slug", "number")
		})

	query = query.Scopes(r.WithPaginate(paging))

	if err := query.Find(&comments).Error; err != nil {
		return nil, 0, err
	}

	return comments, total, nil
}
