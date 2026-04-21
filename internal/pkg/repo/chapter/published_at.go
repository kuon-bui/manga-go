package chapterrepo

import (
	"database/sql"
	"time"

	"manga-go/internal/pkg/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (r *ChapterRepo) GetLatestPublishedChapterTimeByComicIDWithTransaction(tx *gorm.DB, comicID uuid.UUID) (*time.Time, error) {
	var latestPublishedAt sql.NullTime

	err := tx.Model(&model.Chapter{}).
		Select("MAX(published_at)").
		Where("comic_id = ? AND is_published = ? AND deleted_at IS NULL", comicID, true).
		Scan(&latestPublishedAt).Error
	if err != nil {
		return nil, err
	}

	if !latestPublishedAt.Valid {
		return nil, nil
	}

	publishedAt := latestPublishedAt.Time
	return &publishedAt, nil
}
