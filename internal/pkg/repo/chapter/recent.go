package chapterrepo

import (
	"context"
	"manga-go/internal/pkg/common"
	"manga-go/internal/pkg/model"
)

func (r *ChapterRepo) FindRecentUpdates(ctx context.Context, paging *common.Paging) ([]*model.Chapter, int64, error) {
	var chapters []*model.Chapter
	var total int64

	baseQuery := r.DB.WithContext(ctx).
		Model(&model.Chapter{}).
		Where("is_published = ?", true)

	if err := baseQuery.
		Distinct("comic_id").
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	latestChapterIDs := baseQuery.
		Select("DISTINCT ON (comic_id) id").
		Order("comic_id, published_at DESC NULLS LAST, created_at DESC, id DESC")

	query := r.DB.WithContext(ctx).
		Model(&model.Chapter{}).
		Where("chapters.id IN (?)", latestChapterIDs).
		Order("chapters.published_at DESC NULLS LAST, chapters.created_at DESC").
		Preload("Comic")

	query = query.Scopes(r.WithPaginate(paging))

	if err := query.Find(&chapters).Error; err != nil {
		return nil, 0, err
	}

	return chapters, total, nil
}
