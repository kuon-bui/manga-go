package roleservice

import (
	"context"
	"errors"
	"manga-go/internal/app/api/common/response"
	"manga-go/internal/pkg/authorization"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (s *RoleService) RemovePermission(ctx context.Context, roleID, permissionID uuid.UUID) response.Result {
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

	perm, err := s.permissionRepo.FindOne(ctx, []any{
		clause.Eq{Column: "id", Value: permissionID},
	}, nil)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.ResultNotFound("Permission")
		}
		s.logger.Error("Failed to find permission", "error", err)
		return response.ResultErrDb(err)
	}

	if err := s.roleRepo.RemovePermission(ctx, roleID, perm); err != nil {
		s.logger.Error("Failed to remove permission from role", "error", err)
		return response.ResultErrDb(err)
	}

	if s.policyManager != nil {
		if err := s.policyManager.RemovePermissionForRole(role.Name, perm.Name, authorization.OrgPlatform); err != nil {
			s.logger.Error("Failed to remove authorization policy", "error", err)
			return response.ResultErrInternal(err)
		}
	}

	return response.ResultSuccess("Permission removed successfully", nil)
}
