package chapterserivce

import (
	"context"
	"errors"
	"manga-go/internal/app/api/common/response"
	"manga-go/internal/pkg/common"
	"manga-go/internal/pkg/model"
	"manga-go/internal/pkg/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (s *ChapterService) GetChapter(ctx context.Context, chapterSlug string) response.Result {
	comicID, ok := common.GetComicIdFromContext(ctx)
	if !ok {
		return response.ResultError("Comic not found in context")
	}

	chapter, err := s.chapterRepo.FindOne(ctx, []any{
		clause.Eq{Column: "comic_id", Value: comicID},
		clause.Eq{Column: "slug", Value: chapterSlug},
	}, map[string]common.MoreKeyOption{
		"Pages": {},
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.ResultNotFound("Chapter")
		}
		s.logger.Error("Failed to find chapter", "error", err)
		return response.ResultErrDb(err)
	}

	if err := s.enrichChapterPageReactions(ctx, chapter); err != nil {
		s.logger.Error("Failed to enrich chapter page reactions", "error", err)
		return response.ResultErrDb(err)
	}

	return response.ResultSuccess("Chapter retrieved successfully", chapter)
}

func (s *ChapterService) enrichChapterPageReactions(ctx context.Context, chapter *model.Chapter) error {
	if chapter == nil || len(chapter.Pages) == 0 {
		return nil
	}

	pageIds := make([]uuid.UUID, 0, len(chapter.Pages))
	for _, page := range chapter.Pages {
		if page != nil {
			pageIds = append(pageIds, page.ID)
		}
	}
	if len(pageIds) == 0 {
		return nil
	}

	reactionCounts, err := s.pageReactionRepo.CountByPageIds(ctx, pageIds)
	if err != nil {
		return err
	}

	userReactions := make(map[uuid.UUID]string)
	user, err := utils.GetCurrentUserFormContext(ctx)
	if err == nil && user != nil {
		userReactions, err = s.pageReactionRepo.GetUserReactionsByPageIds(ctx, pageIds, user.ID)
		if err != nil {
			return err
		}
	}

	for _, page := range chapter.Pages {
		if page == nil {
			continue
		}

		counts := model.ReactionCounts{}
		if countMap, ok := reactionCounts[page.ID]; ok {
			counts.LIKE = countMap[common.ReactionTypeLike]
			counts.LOVE = countMap[common.ReactionTypeLove]
			counts.HAHA = countMap[common.ReactionTypeHaha]
			counts.WOW = countMap[common.ReactionTypeWow]
			counts.SAD = countMap[common.ReactionTypeSad]
			counts.ANGRY = countMap[common.ReactionTypeAngry]
		}
		page.ReactionCounts = counts

		if reaction, ok := userReactions[page.ID]; ok && reaction != "" {
			reactionValue := reaction
			page.UserReaction = &reactionValue
		} else {
			page.UserReaction = nil
		}
	}

	return nil
}
