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

func (s *TranslationGroupService) GrantPermission(ctx context.Context, ownerID uuid.UUID, slug string, req *translationgrouprequest.GrantPermissionRequest) response.Result {
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

	// Verify the requester is the group owner
	if group.OwnerID != ownerID {
		return response.ResultForbidden()
	}

	// Verify the target user is a member of this group
	_, err = s.userRepo.FindOne(ctx, []any{
		clause.Eq{Column: "id", Value: req.MemberID},
		clause.Eq{Column: "translation_group_id", Value: group.ID},
	}, nil)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.ResultError("User is not a member of this group")
		}
		s.logger.Error("Failed to find member", "error", err)
		return response.ResultErrDb(err)
	}

	groupIDStr := group.ID.String()
	memberIDStr := req.MemberID.String()

	// Grant chapter_creator role in the group's Casbin domain
	if _, err := s.enforcer.AddRoleForUserInDomain(memberIDStr, "chapter_creator", groupIDStr); err != nil {
		s.logger.Errorf("Failed to grant chapter_creator role to user %s in group %s: %v", memberIDStr, groupIDStr, err)
		return response.ResultErrInternal(err)
	}

	return response.ResultSuccess("Chapter creation permission granted successfully", nil)
}
