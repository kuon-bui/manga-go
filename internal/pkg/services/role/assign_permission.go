package roleservice

import (
	"context"
	"errors"
	"manga-go/internal/app/api/common/response"
	"manga-go/internal/pkg/authorization"
	rolerequest "manga-go/internal/pkg/request/role"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// AssignPermissions replaces the role's grants with exactly the names given.
func (s *RoleService) AssignPermissions(
	ctx context.Context,
	roleID uuid.UUID,
	req *rolerequest.AssignPermissionRequest,
	expectedVersion ...string,
) response.Result {
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

	// Validate the whole batch before writing: a single unknown name must not
	// leave the role stripped of the grants it already had.
	for _, name := range *req.Permissions {
		if err := authorization.ValidatePermissionName(name); err != nil {
			return response.ResultError(err.Error())
		}
	}
	if s.authAdmin != nil && s.authAdmin.MutationReady() {
		return s.authAdmin.ReplaceRolePermissions(ctx, roleID, *req.Permissions, roleVersion(expectedVersion))
	}

	if err := s.policyManager.ReplacePermissionsForRole(
		role.ID.String(), *req.Permissions, authorization.OrgPlatform,
	); err != nil {
		s.logger.Error("Failed to update authorization policy", "error", err)
		return response.ResultErrInternal(err)
	}

	return response.ResultSuccess("Permissions assigned successfully", *req.Permissions)
}

func roleVersion(values []string) string {
	if len(values) == 0 {
		return ""
	}
	return values[0]
}
