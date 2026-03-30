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

func (s *TranslationGroupService) JoinTranslationGroup(ctx context.Context, userID uuid.UUID, req *translationgrouprequest.JoinTranslationGroupRequest) response.Result {
	group, err := s.translationGroupRepo.FindOne(ctx, []any{
		clause.Eq{Column: "slug", Value: req.Slug},
	}, nil)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.ResultNotFound("TranslationGroup")
		}
		s.logger.Error("Failed to find translation group", "error", err)
		return response.ResultErrDb(err)
	}

	// Check that user is not already in a group
	user, err := s.userRepo.FindOne(ctx, []any{
		clause.Eq{Column: "id", Value: userID},
	}, nil)
	if err != nil {
		s.logger.Error("Failed to find user", "error", err)
		return response.ResultErrDb(err)
	}

	if user.TranslationGroupID != nil {
		return response.ResultError("User is already a member of a translation group")
	}

	// Set the user's translation_group_id
	if err := s.userRepo.Update(ctx, []any{
		clause.Eq{Column: "id", Value: userID},
	}, map[string]any{
		"translation_group_id": group.ID,
	}); err != nil {
		s.logger.Error("Failed to join translation group", "error", err)
		return response.ResultErrDb(err)
	}

	groupIDStr := group.ID.String()

	// Assign group_member role in the group's Casbin domain
	if _, err := s.enforcer.AddRoleForUserInDomain(userID.String(), "group_member", groupIDStr); err != nil {
		s.logger.Errorf("Failed to assign group_member role to user %s for group %s: %v", userID, groupIDStr, err)
		return response.ResultErrInternal(err)
	}

	return response.ResultSuccess("Joined translation group successfully", group)
}
