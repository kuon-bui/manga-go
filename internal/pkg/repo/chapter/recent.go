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
		Joins("JOIN comics ON comics.id = chapters.comic_id").
		Where("chapters.is_published = ? AND comics.is_published = ?", true, true)

	if err := baseQuery.
		Distinct("comic_id").
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	latestChapterIDs := r.DB.WithContext(ctx).
		Table("(SELECT DISTINCT ON (chapters.comic_id) chapters.id FROM chapters JOIN comics ON comics.id = chapters.comic_id WHERE chapters.is_published = true AND comics.is_published = true ORDER BY chapters.comic_id, chapters.published_at DESC NULLS LAST, chapters.created_at DESC, chapters.id DESC) AS latest").
		Select("id")

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
