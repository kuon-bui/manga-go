package tagrepo

import (
	"context"
	"manga-go/internal/pkg/model"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

const slugToIdCacheKey = "slug:tag"

func (r *TagRepo) GetIdBySlug(ctx context.Context, slug string) (id uuid.UUID, err error) {
	// get from cache
	idStr := ""
	r.rds.Client().HGet(ctx, slugToIdCacheKey, slug).Scan(&idStr)
	if idStr != "" {
		id, err = uuid.Parse(idStr)
		if err == nil {
			return id, nil
		}
	}

	var tag model.Tag
	if err := r.DB.WithContext(ctx).Select("id").Where("slug = ?", slug).First(&tag).Error; err != nil {
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
		tag.ID.String(),
	)

	return tag.ID, nil
}
