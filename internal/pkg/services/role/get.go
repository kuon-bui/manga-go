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

func (s *RoleService) GetRole(ctx context.Context, id uuid.UUID) response.Result {
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

	// Grants live in the policy engine, not alongside the role metadata.
	permissions, err := s.policyManager.PermissionNamesForRole(role.ID.String(), authorization.OrgPlatform)
	if err != nil {
		s.logger.Error("Failed to read role permissions", "error", err)
		return response.ResultErrInternal(err)
	}
	role.Permissions = permissions

	return response.ResultSuccess("Role retrieved successfully", role)
}
