package chapterrepo

import (
	"context"
	"manga-go/internal/pkg/model"

	"github.com/google/uuid"
)

func (r *ChapterRepo) GetNextChapterIdx(ctx context.Context, comicID uuid.UUID) (uint, error) {
	var maxChapterIdx int64
	err := r.DB.WithContext(ctx).
		Model(&model.Chapter{}).
		Select("COALESCE(MAX(chapter_idx), -1)").
		Where("comic_id = ?", comicID).
		Scan(&maxChapterIdx).Error
	if err != nil {
		return 0, err
	}

	return uint(maxChapterIdx + 1), nil
}
