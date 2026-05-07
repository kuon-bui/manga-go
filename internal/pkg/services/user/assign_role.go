package userservice

import (
	"context"
	"errors"
	"manga-go/internal/app/api/common/response"
	"manga-go/internal/pkg/authorization"
	"manga-go/internal/pkg/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (s *UserService) AssignRoles(ctx context.Context, userID uuid.UUID, roleIDs []uuid.UUID) response.Result {
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

	roles := make([]*model.Role, 0, len(roleIDs))
	for _, id := range roleIDs {
		role, err := s.roleRepo.FindOne(ctx, []any{
			clause.Eq{Column: "id", Value: id},
		}, nil)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return response.ResultNotFound("Role")
			}
			s.logger.Error("Failed to find role", "error", err)
			return response.ResultErrDb(err)
		}
		roles = append(roles, role)
	}

	if err := s.userRepo.AssignRoles(ctx, userID, roles); err != nil {
		s.logger.Error("Failed to assign roles to user", "error", err)
		return response.ResultErrDb(err)
	}

	if s.policyManager != nil {
		roleNames := make([]string, 0, len(roles))
		for _, role := range roles {
			roleNames = append(roleNames, role.Name)
		}
		if err := s.policyManager.ReplaceRolesForUser(userID.String(), roleNames, authorization.OrgPlatform); err != nil {
			s.logger.Error("Failed to update authorization policy", "error", err)
			return response.ResultErrInternal(err)
		}
	}

	return response.ResultSuccess("Roles assigned successfully", nil)
}
