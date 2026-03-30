package comicrepo

import (
	"context"
	"manga-go/internal/pkg/model"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

const slugToIdCacheKey = "slug:comic"

func (r *ComicRepo) GetIdBySlug(ctx context.Context, slug string) (id uuid.UUID, err error) {
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
	var comic model.Comic
	err = r.DB.WithContext(ctx).Select("id").Where("slug = ?", slug).First(&comic).Error
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
		comic.ID.String(),
	)

	return comic.ID, nil
}

// GetIdAndGroupIdBySlug returns the comic ID and its associated translation group ID for a given slug.
func (r *ComicRepo) GetIdAndGroupIdBySlug(ctx context.Context, slug string) (uuid.UUID, *uuid.UUID, error) {
	var comic model.Comic
	err := r.DB.WithContext(ctx).Select("id", "translation_group_id").Where("slug = ?", slug).First(&comic).Error
	if err != nil {
		return uuid.Nil, nil, err
	}
	return comic.ID, comic.TranslationGroupID, nil
}
