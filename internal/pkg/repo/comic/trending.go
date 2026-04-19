package comicrepo

import (
	"context"
	"manga-go/internal/pkg/common"
	"manga-go/internal/pkg/model"

	"gorm.io/gorm"
)

func (r *ComicRepo) FindTrending(ctx context.Context, limit int) ([]*model.Comic, error) {
	var comics []*model.Comic

	db := r.DB.WithContext(ctx).
		Model(&model.Comic{}).
		Where("is_published = ?", true).
		Select(statsSelect).
		Order("follow_count DESC NULLS LAST").
		Limit(limit)

	db = r.ApplyPreloadMoreKeys(db, map[string]common.MoreKeyOption{
		"Authors": {},
		"Genres":  {},
		"Chapters": {
			Custom: func(tx *gorm.DB) *gorm.DB {
				return tx.Order("chapter_idx DESC").Limit(1)
			},
		},
	})

	if err := db.Find(&comics).Error; err != nil {
		return nil, err
	}

	return comics, nil
}
