package translationgroupservice

import (
	"context"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm/clause"
)

const slugToIdCacheKey = "slug:translation_group"

func (s *TranslationGroupService) GetTranslationGroupIDBySlug(ctx context.Context, slug string) (uuid.UUID, error) {
	// get from cache
	idStr := ""
	s.rds.Client().HGet(ctx, slugToIdCacheKey, slug).Scan(&idStr)
	if idStr != "" {
		id, err := uuid.Parse(idStr)
		if err == nil {
			return id, nil
		}
	}

	// get from db
	translationGroup, err := s.translationGroupRepo.FindOne(ctx, []any{
		clause.Eq{Column: "slug", Value: slug},
	}, nil)
	if err != nil {
		return uuid.Nil, err
	}

	// set to cache
	s.rds.Client().HSetEXWithArgs(
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
