package translationgrouprepo

import (
	"context"
	"manga-go/internal/pkg/model"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

const slugToIdCacheKey = "slug:translation_group"

func (r *TranslationGroupRepo) GetIdBySlug(ctx context.Context, slug string) (uuid.UUID, error) {
	// get from cache
	idStr := ""
	r.rds.Client().HGet(ctx, slugToIdCacheKey, slug).Scan(&idStr)
	if idStr != "" {
		id, err := uuid.Parse(idStr)
		if err == nil {
			return id, nil
		}
	}

	// get from db
	var translationGroup model.TranslationGroup
	err := r.DB.WithContext(ctx).Select("id").Where("slug = ?", slug).First(&translationGroup).Error
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
		translationGroup.ID.String(),
	)

	return translationGroup.ID, nil
}
