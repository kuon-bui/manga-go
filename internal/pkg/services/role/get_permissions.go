package roleservice

import (
	"context"
	"errors"
	"manga-go/internal/app/api/common/response"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (s *RoleService) GetRolePermissions(ctx context.Context, roleID uuid.UUID) response.Result {
	_, err := s.roleRepo.FindOne(ctx, []any{
		clause.Eq{Column: "id", Value: roleID},
	}, nil)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.ResultNotFound("Role")
		}
		s.logger.Error("Failed to find role", "error", err)
		return response.ResultErrDb(err)
	}

	permissions, err := s.roleRepo.GetPermissions(ctx, roleID)
	if err != nil {
		s.logger.Error("Failed to get role permissions", "error", err)
		return response.ResultErrDb(err)
	}

	return response.ResultSuccess("Role permissions retrieved successfully", permissions)
}
