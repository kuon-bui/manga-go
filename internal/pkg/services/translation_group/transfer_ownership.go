package translationgroupservice

import (
	"context"
	"errors"
	"manga-go/internal/app/api/common/response"
	translationgrouprequest "manga-go/internal/pkg/request/translation_group"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (s *TranslationGroupService) TransferOwnership(ctx context.Context, slug string, req *translationgrouprequest.TransferOwnershipRequest) response.Result {
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

	// Verify the new owner is a member of this group
	_, err = s.userRepo.FindOne(ctx, []any{
		clause.Eq{Column: "id", Value: req.NewOwnerID},
		clause.Eq{Column: "translation_group_id", Value: group.ID},
	}, nil)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.ResultError("New owner must be a member of the translation group")
		}
		s.logger.Error("Failed to find new owner", "error", err)
		return response.ResultErrDb(err)
	}

	if err := s.translationGroupRepo.Update(ctx, []any{
		clause.Eq{Column: "id", Value: group.ID},
	}, map[string]any{
		"owner_id": req.NewOwnerID,
	}); err != nil {
		s.logger.Error("Failed to transfer ownership", "error", err)
		return response.ResultErrDb(err)
	}

	group.OwnerID = req.NewOwnerID
	return response.ResultSuccess("Ownership transferred successfully", group)
}
