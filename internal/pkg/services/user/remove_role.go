package userservice

import (
	"context"
	"errors"
	"manga-go/internal/app/api/common/response"
	"manga-go/internal/pkg/authorization"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (s *UserService) RemoveRole(
	ctx context.Context,
	userID, roleID uuid.UUID,
	expectedVersion ...string,
) response.Result {
	_, err := s.userRepo.FindOne(ctx, []any{
		clause.Eq{Column: "id", Value: userID},
	}, nil)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.ResultNotFound("User")
		}
		s.logger.Error("Failed to find user", "error", err)
		return response.ResultErrDb(err)
	}

	role, err := s.roleRepo.FindOne(ctx, []any{
		clause.Eq{Column: "id", Value: roleID},
	}, nil)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.ResultNotFound("Role")
		}
		s.logger.Error("Failed to find role", "error", err)
		return response.ResultErrDb(err)
	}

	if s.authAdmin != nil && s.authAdmin.MutationReady() {
		current, err := s.policyManager.RolesForUser(userID.String(), authorization.OrgPlatform)
		if err != nil {
			return response.ResultErrInternal(err)
		}
		remaining := make([]uuid.UUID, 0, len(current))
		for _, rawID := range current {
			id, err := uuid.Parse(rawID)
			if err != nil {
				return response.ResultErrInternal(err)
			}
			if id != roleID {
				remaining = append(remaining, id)
			}
		}
		return s.authAdmin.ReplaceUserRoles(ctx, userID, remaining, firstVersion(expectedVersion))
	}

	if err := s.policyManager.RemoveRoleForUser(userID.String(), role.ID.String(), authorization.OrgPlatform); err != nil {
		s.logger.Error("Failed to update authorization policy", "error", err)
		return response.ResultErrInternal(err)
	}

	return response.ResultSuccess("Role removed successfully", nil)
}
