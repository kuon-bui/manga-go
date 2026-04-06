package chapterserivce

import (
	"context"
	"errors"
	"fmt"
	"manga-go/internal/app/api/common/response"
	"manga-go/internal/pkg/bitset"
	"manga-go/internal/pkg/model"
	"manga-go/internal/pkg/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (s *ChapterService) MarkChapterAsRead(ctx context.Context, comicID, chapterID uuid.UUID) response.Result {
	user, err := utils.GetCurrentUserFormContext(ctx)
	if err != nil {
		s.logger.Error("Failed to get current user from context", "error", err)
		return response.ResultErrInternal(err)
	}

	// Check if the chapter exists
	chapter, err := s.chapterRepo.FindOne(ctx, []any{
		clause.Eq{Column: "id", Value: chapterID},
		clause.Eq{Column: "comic_id", Value: comicID},
	}, nil)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.ResultNotFound("Chapter")
		}
		s.logger.Error("Failed to find chapter", "error", err)
		return response.ResultErrDb(err)
	}

	// Update or create the UserComicRead record
	var userComicRead *model.UserComicRead
	userComicRead, err = s.userComicReadRepo.FindOne(ctx, []any{
		clause.Eq{Column: "user_id", Value: user.ID},
		clause.Eq{Column: "comic_id", Value: comicID},
	}, nil)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			s.logger.Error("Failed to find user comic read data", "error", err)
			return response.ResultErrDb(err)
		}

		// Create a new record if not found
		count, err := s.chapterRepo.CountAll(ctx, []any{clause.Eq{Column: "comic_id", Value: comicID}})
		if err != nil {
			s.logger.Error("Failed to count chapters for comic", "error", err)
			return response.ResultErrInternal(err)
		}
		// Mark all chapters up to the current one as read
		readData := bitset.NewReadBitset(int(count))
		userComicRead = &model.UserComicRead{
			UserID:   user.ID,
			ComicID:  comicID,
			ReadData: readData,
		}
	}

	fmt.Printf("UserComicRead before marking as read: %+v\n", *userComicRead.ReadData)
	// Mark the chapter as read
	userComicRead.ReadData.Mark(int(chapter.ChapterIdx))
	err = s.userComicReadRepo.Save(ctx, userComicRead)
	if err != nil {
		s.logger.Error("Failed to save user comic read data", "error", err)
		return response.ResultErrInternal(err)
	}

	return response.ResultSuccess("Chapter marked as read successfully", nil)
}
