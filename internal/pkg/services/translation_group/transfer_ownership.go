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

func (s *TranslationGroupService) TransferOwnership(ctx context.Context, requesterID uuid.UUID, slug string, req *translationgrouprequest.TransferOwnershipRequest) response.Result {
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

	groupIDStr := group.ID.String()

	// Update Casbin roles: assign group_owner to new owner, demote old owner to chapter_creator
	if _, err := s.enforcer.AddRoleForUserInDomain(req.NewOwnerID.String(), "group_owner", groupIDStr); err != nil {
		s.logger.Errorf("Failed to assign group_owner to new owner %s in group %s: %v", req.NewOwnerID, groupIDStr, err)
	}
	if _, err := s.enforcer.DeleteRoleForUserInDomain(requesterID.String(), "group_owner", groupIDStr); err != nil {
		s.logger.Errorf("Failed to remove group_owner from old owner %s in group %s: %v", requesterID, groupIDStr, err)
	}

	group.OwnerID = req.NewOwnerID
	return response.ResultSuccess("Ownership transferred successfully", group)
}
