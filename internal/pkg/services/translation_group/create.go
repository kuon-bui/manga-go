package translationgroupservice

import (
	"context"
	"manga-go/internal/app/api/common/response"
	"manga-go/internal/pkg/model"
	translationgrouprequest "manga-go/internal/pkg/request/translation_group"

	"github.com/google/uuid"
	"gorm.io/gorm/clause"
)

func (s *TranslationGroupService) CreateTranslationGroup(ctx context.Context, ownerID uuid.UUID, req *translationgrouprequest.CreateTranslationGroupRequest) response.Result {
	group := model.TranslationGroup{
		Name:    req.Name,
		Slug:    req.Slug,
		OwnerID: ownerID,
	}

	if err := s.translationGroupRepo.Create(ctx, &group); err != nil {
		s.logger.Error("Failed to create translation group", "error", err)
		return response.ResultErrDb(err)
	}

	// Set the creator's translation_group_id
	if err := s.userRepo.Update(ctx, []any{
		clause.Eq{Column: "id", Value: ownerID},
	}, map[string]any{
		"translation_group_id": group.ID,
	}); err != nil {
		s.logger.Error("Failed to set user translation group", "error", err)
		return response.ResultErrDb(err)
	}

	return response.ResultSuccess("Translation group created successfully", group)
}
