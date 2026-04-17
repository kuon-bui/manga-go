package chapterserivce

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm/clause"
)

const slugToIdCacheKey = "slug:chapter"
const comicHasChapterKey = "comic:has_chapter"

func (s *ChapterService) GetChapterIDBySlug(ctx context.Context, comicID uuid.UUID, slug string) (uuid.UUID, error) {
	defaultErr := errors.New("chapter not found")
	idStr := ""
	s.rds.Client().HGet(ctx, slugToIdCacheKey, slug).Scan(&idStr)
	if idStr != "" {
		id, err := uuid.Parse(idStr)
		if err == nil {
			return id, nil
		}
	}

	// if cache has record, it means the chapter does not exist
	hasInRedis, _ := s.rds.Client().HExists(ctx, comicHasChapterKey, comicID.String()).Result()
	if hasInRedis {
		var comicHasChapter string
		s.rds.Client().HGet(ctx, comicHasChapterKey, comicID.String()).Scan(&comicHasChapter)
		// if cache has record but not "1", it means the chapter does not exist
		if comicHasChapter != "1" {
			return uuid.Nil, defaultErr
		}
	}
	chapter, err := s.chapterRepo.FindOne(ctx, []any{
		clause.Eq{Column: "comic_id", Value: comicID},
		clause.Eq{Column: "slug", Value: slug},
	}, nil)
	if err != nil {
		return uuid.Nil, err
	}

	// set chapter id to cache
	s.rds.Client().HSetEXWithArgs(
		ctx,
		slugToIdCacheKey,
		&redis.HSetEXOptions{
			ExpirationType: redis.HSetEXExpirationEX,
			ExpirationVal:  10 * 60, // 10 minutes
		},
		slug,
		chapter.ID.String(),
	)

	// set comic has chapter to cache
	s.rds.Client().HSetEXWithArgs(
		ctx,
		comicHasChapterKey,
		&redis.HSetEXOptions{
			ExpirationType: redis.HSetEXExpirationEX,
			ExpirationVal:  10 * 60, // 10 minutes
		},
		comicID.String(),
		"1",
	)

	return chapter.ID, nil
}
