package permissionservice

import (
	"context"
	"errors"
	"manga-go/internal/app/api/common/response"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (s *PermissionService) GetPermission(ctx context.Context, id uuid.UUID) response.Result {
	permission, err := s.permissionRepo.FindOne(ctx, []any{
		clause.Eq{Column: "id", Value: id},
	}, nil)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.ResultNotFound("Permission")
		}
		s.logger.Error("Failed to find permission", "error", err)
		return response.ResultErrDb(err)
	}

	return response.ResultSuccess("Permission retrieved successfully", permission)
}
