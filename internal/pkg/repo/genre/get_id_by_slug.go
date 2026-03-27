package genrerepo

import (
	"context"
	"manga-go/internal/pkg/model"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

const slugToIdCacheKey = "slug:genre"

func (r *GenreRepo) GetIdBySlug(ctx context.Context, slug string) (id uuid.UUID, err error) {
	// get from db
	idStr := ""
	r.rds.Client().HGet(ctx, slugToIdCacheKey, slug).Scan(&idStr)
	if idStr != "" {
		id, err = uuid.Parse(idStr)
		if err == nil {
			return id, nil
		}
	}

	var genre model.Genre
	if err := r.DB.WithContext(ctx).Select("id").Where("slug = ?", slug).First(&genre).Error; err != nil {
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
		genre.ID.String(),
	)

	return genre.ID, nil
}
