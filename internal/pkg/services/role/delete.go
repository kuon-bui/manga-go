package roleservice

import (
	"context"
	"errors"
	"manga-go/internal/app/api/common/response"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (s *RoleService) DeleteRole(ctx context.Context, id uuid.UUID) response.Result {
	_, err := s.roleRepo.FindOne(ctx, []any{
		clause.Eq{Column: "id", Value: id},
	}, nil)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.ResultNotFound("Role")
		}
		s.logger.Error("Failed to find role for deletion", "error", err)
		return response.ResultErrDb(err)
	}

	if err := s.roleRepo.DeleteSoft(ctx, []any{
		clause.Eq{Column: "id", Value: id},
	}); err != nil {
		s.logger.Error("Failed to delete role", "error", err)
		return response.ResultErrDb(err)
	}

	return response.ResultSuccess("Role deleted successfully", nil)
}
