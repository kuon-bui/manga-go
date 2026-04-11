package comicservice

import (
	"context"
	"errors"
	"manga-go/internal/app/api/common/response"
	"manga-go/internal/pkg/common"
	"manga-go/internal/pkg/utils"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (s *ComicService) GetComic(ctx context.Context, slug string) response.Result {
	comic, err := s.comicRepo.FindOne(ctx, []any{
		clause.Eq{Column: "slug", Value: slug},
	}, map[string]common.MoreKeyOption{
		"Authors":  {},
		"Artists":  {},
		"Genres":   {},
		"Tags":     {},
		"Chapters": {},
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.ResultNotFound("Comic")
		}
		s.logger.Error("Failed to find comic", "error", err)
		return response.ResultErrDb(err)
	}

	user, err := utils.GetCurrentUserFormContext(ctx)
	if err != nil {
		s.logger.Error("Failed to get current user from context", "error", err)
		return response.ResultErrInternal(err)
	}

	userComicRead, err := s.userComicReadRepo.FindOne(ctx, []any{
		clause.Eq{Column: "user_id", Value: user.ID},
		clause.Eq{Column: "comic_id", Value: comic.ID},
	}, nil)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		s.logger.Error("Failed to find user comic read data", "error", err)
		return response.ResultErrDb(err)
	}

	if userComicRead != nil {
		for _, chapter := range comic.Chapters {
			if userComicRead.ReadData.IsRead(int(chapter.ChapterIdx)) {
				chapter.IsRead = true
			}
		}
	}

	return response.ResultSuccess("Comic retrieved successfully", comic)
}
