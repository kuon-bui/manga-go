package pageservice

import (
	"context"
	"errors"
	"manga-go/internal/app/api/common/response"
	"manga-go/internal/pkg/common"
	"manga-go/internal/pkg/model"
	pagerequest "manga-go/internal/pkg/request/page"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (s *PageService) HandleReaction(ctx context.Context, user *model.User, pageId uuid.UUID, req *pagerequest.AddReactionRequest) response.Result {
	reactionType := common.NormalizeReactionType(req.Type)
	if !common.IsValidReactionType(reactionType) {
		return response.ResultError("Invalid reaction type")
	}

	page, err := s.pageRepo.FindOne(ctx, []any{
		clause.Eq{Column: "id", Value: pageId},
	}, nil)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.ResultNotFound("Page")
		}
		s.logger.Error("Failed to find page", "error", err)
		return response.ResultErrDb(err)
	}

	hadReacted, err := s.pageReactionRepo.ExistsByPageIdAndUserId(ctx, pageId, user.ID)
	if err != nil {
		s.logger.Error("Failed to check page reaction existence", "error", err)
		return response.ResultErrDb(err)
	}

	if hadReacted {
		return s.removeReaction(ctx, user, page)
	}

	return s.addReaction(ctx, user, page, reactionType)
}

func (s *PageService) addReaction(ctx context.Context, user *model.User, page *model.Page, reactionType string) response.Result {
	reaction := &model.PageReaction{
		PageId: page.ID,
		Reaction: model.Reaction{
			UserId: user.ID,
			Type:   reactionType,
		},
	}

	if err := s.pageReactionRepo.Create(ctx, reaction); err != nil {
		s.logger.Error("Failed to add page reaction", "error", err)
		return response.ResultErrDb(err)
	}

	return response.ResultSuccess("Page reaction added successfully", reaction)
}

func (s *PageService) removeReaction(ctx context.Context, user *model.User, page *model.Page) response.Result {
	if err := s.pageReactionRepo.DeleteSoft(ctx, []any{
		clause.Eq{Column: "page_id", Value: page.ID},
		clause.Eq{Column: "user_id", Value: user.ID},
	}); err != nil {
		s.logger.Error("Failed to remove page reaction", "error", err)
		return response.ResultErrDb(err)
	}

	return response.ResultSuccess("Page reaction removed successfully", nil)
}
