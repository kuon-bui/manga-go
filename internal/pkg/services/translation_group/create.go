package translationgroupservice

import (
	"context"
	"manga-go/internal/app/api/common/response"
	casbinpkg "manga-go/internal/pkg/casbin"
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

	groupIDStr := group.ID.String()

	// Seed per-group Casbin policies
	casbinpkg.SeedGroupPolicies(s.enforcer, groupIDStr, s.logger)

	// Assign group_owner role to the creator within this group's domain
	if _, err := s.enforcer.AddRoleForUserInDomain(ownerID.String(), "group_owner", groupIDStr); err != nil {
		s.logger.Errorf("Failed to assign group_owner role to user %s for group %s: %v", ownerID, groupIDStr, err)
		return response.ResultErrInternal(err)
	}

	return response.ResultSuccess("Translation group created successfully", group)
}
