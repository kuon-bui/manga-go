package chapterrepo

import (
	"context"
	"manga-go/internal/pkg/model"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

const slugToIdCacheKey = "slug:chapter"

func (r *ChapterRepo) GetIdBySlug(ctx context.Context, slug string) (id uuid.UUID, err error) {
	// get from cache
	idStr := ""
	r.rds.Client().HGet(ctx, slugToIdCacheKey, slug).Scan(&idStr)
	if idStr != "" {
		id, err = uuid.Parse(idStr)
		if err == nil {
			return id, nil
		}
	}

	// get from db
	var chapter model.Chapter
	err = r.DB.WithContext(ctx).Select("id").Where("slug = ?", slug).First(&chapter).Error
	if err != nil {
		return uuid.Nil, err
	}

	// set to cache
	r.rds.Client().HSetEXWithArgs(
		ctx,
		slugToIdCacheKey,
		&redis.HSetEXOptions{
			ExpirationType: redis.HSetEXExpirationEX,
			ExpirationVal:  10 * 60, // 10 minutes
		},
		slug,
		chapter.ID.String(),
	)

	return chapter.ID, nil
}
