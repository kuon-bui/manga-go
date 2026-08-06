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

func (s *RoleService) RemovePermission(ctx context.Context, roleID uuid.UUID, name string) response.Result {
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

	if err := authorization.ValidatePermissionName(name); err != nil {
		return response.ResultError(err.Error())
	}

	if err := s.policyManager.RevokePermissionFromRole(
		role.ID.String(), name, authorization.OrgPlatform,
	); err != nil {
		s.logger.Error("Failed to update authorization policy", "error", err)
		return response.ResultErrInternal(err)
	}

	return response.ResultSuccess("Permission removed successfully", nil)
}
