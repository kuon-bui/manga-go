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

func (s *TranslationGroupService) KickMember(ctx context.Context, ownerID uuid.UUID, slug string, req *translationgrouprequest.KickMemberRequest) response.Result {
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

	// Prevent kicking yourself (the owner)
	if req.MemberID == ownerID {
		return response.ResultError("Group owner cannot kick themselves")
	}

	// Verify the target user is actually in this group
	member, err := s.userRepo.FindOne(ctx, []any{
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

	// Remove the user's translation_group_id
	if err := s.userRepo.Update(ctx, []any{
		clause.Eq{Column: "id", Value: member.ID},
	}, map[string]any{
		"translation_group_id": nil,
	}); err != nil {
		s.logger.Error("Failed to kick member", "error", err)
		return response.ResultErrDb(err)
	}

	groupIDStr := group.ID.String()
	memberIDStr := req.MemberID.String()

	// Remove all Casbin roles for this user in the group's domain
	if _, err := s.enforcer.DeleteRoleForUserInDomain(memberIDStr, "group_member", groupIDStr); err != nil {
		s.logger.Errorf("Failed to remove group_member role for user %s in group %s: %v", memberIDStr, groupIDStr, err)
		return response.ResultErrInternal(err)
	}
	// chapter_creator may or may not exist – log but don't fail on this removal
	if _, err := s.enforcer.DeleteRoleForUserInDomain(memberIDStr, "chapter_creator", groupIDStr); err != nil {
		s.logger.Errorf("Failed to remove chapter_creator role for user %s in group %s: %v", memberIDStr, groupIDStr, err)
	}

	return response.ResultSuccess("Member kicked successfully", nil)
}
