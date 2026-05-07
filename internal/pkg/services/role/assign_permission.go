package roleservice

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

func (s *RoleService) AssignPermissions(ctx context.Context, roleID uuid.UUID, permissionIDs []uuid.UUID) response.Result {
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

	perms := make([]*model.Permission, 0, len(permissionIDs))
	for _, id := range permissionIDs {
		perm, err := s.permissionRepo.FindOne(ctx, []any{
			clause.Eq{Column: "id", Value: id},
		}, nil)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return response.ResultNotFound("Permission")
			}
			s.logger.Error("Failed to find permission", "error", err)
			return response.ResultErrDb(err)
		}
		perms = append(perms, perm)
	}

	if err := s.roleRepo.AssignPermissions(ctx, roleID, perms); err != nil {
		s.logger.Error("Failed to assign permissions to role", "error", err)
		return response.ResultErrDb(err)
	}

	if s.policyManager != nil {
		permissionNames := make([]string, 0, len(perms))
		for _, perm := range perms {
			permissionNames = append(permissionNames, perm.Name)
		}
		if err := s.policyManager.ReplacePermissionsForRole(role.Name, permissionNames, authorization.OrgPlatform); err != nil {
			s.logger.Error("Failed to update authorization policy", "error", err)
			return response.ResultErrInternal(err)
		}
	}

	return response.ResultSuccess("Permissions assigned successfully", nil)
}
