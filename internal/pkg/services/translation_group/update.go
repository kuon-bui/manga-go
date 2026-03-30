package translationgroupservice

import (
	"context"
	"errors"
	"manga-go/internal/app/api/common/response"
	translationgrouprequest "manga-go/internal/pkg/request/translation_group"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (s *TranslationGroupService) UpdateTranslationGroup(ctx context.Context, requesterID uuid.UUID, slug string, req *translationgrouprequest.UpdateTranslationGroupRequest) response.Result {
	group, err := s.translationGroupRepo.FindOne(ctx, []any{
		clause.Eq{Column: "slug", Value: slug},
	}, nil)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.ResultNotFound("TranslationGroup")
		}
		s.logger.Error("Failed to find translation group", "error", err)
		return response.ResultErrDb(err)
	}

	if group.OwnerID != requesterID {
		return response.ResultForbidden()
	}

	if err := s.translationGroupRepo.Update(ctx, []any{
		clause.Eq{Column: "id", Value: group.ID},
	}, map[string]any{
		"name": req.Name,
		"slug": req.Slug,
	}); err != nil {
		s.logger.Error("Failed to update translation group", "error", err)
		return response.ResultErrDb(err)
	}

	group.Name = req.Name
	group.Slug = req.Slug
	return response.ResultSuccess("Translation group updated successfully", group)
}
