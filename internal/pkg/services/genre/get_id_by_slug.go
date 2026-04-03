package genreservice

import (
	"context"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm/clause"
)

const slugToIdCacheKey = "slug:genre"

func (s *GenreService) GetGenreIDBySlug(ctx context.Context, slug string) (uuid.UUID, error) {
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
	genre, err := s.genreRepo.FindOne(ctx, []any{
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
		genre.ID.String(),
	)

	return genre.ID, nil
}
