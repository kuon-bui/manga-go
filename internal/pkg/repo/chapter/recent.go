package chapterrepo

import (
	"context"
	"manga-go/internal/pkg/common"
	"manga-go/internal/pkg/model"
)

func (r *ChapterRepo) FindRecentUpdates(ctx context.Context, paging *common.Paging) ([]*model.Chapter, int64, error) {
	var chapters []*model.Chapter
	var total int64

	db := r.DB.WithContext(ctx).
		Model(&model.Chapter{}).
		Where("is_published = ?", true)

	if err := db.
		Distinct("chapters.id").
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	query := r.DB.WithContext(ctx).
		Model(&model.Chapter{}).
		Where("is_published = ?", true).
		Order("chapters.created_at DESC").
		Preload("Comic")

	query = query.Scopes(r.WithPaginate(paging))

	if err := query.Find(&chapters).Error; err != nil {
		return nil, 0, err
	}

	return chapters, total, nil
}
