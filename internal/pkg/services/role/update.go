package roleservice

import (
	"context"
	"errors"
	"manga-go/internal/app/api/common/response"
	rolerequest "manga-go/internal/pkg/request/role"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (s *RoleService) UpdateRole(ctx context.Context, id uuid.UUID, req *rolerequest.UpdateRoleRequest) response.Result {
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

	if err := s.roleRepo.Update(ctx, []any{
		clause.Eq{Column: "id", Value: role.ID},
	}, map[string]any{
		"name": req.Name,
	}); err != nil {
		s.logger.Error("Failed to update role", "error", err)
		return response.ResultErrDb(err)
	}

	role.Name = req.Name
	return response.ResultSuccess("Role updated successfully", role)
}
