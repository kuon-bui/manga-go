package permissionservice

import (
	"context"
	"errors"
	"manga-go/internal/app/api/common/response"
	permissionrequest "manga-go/internal/pkg/request/permission"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (s *PermissionService) UpdatePermission(ctx context.Context, id uuid.UUID, req *permissionrequest.UpdatePermissionRequest) response.Result {
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

	if err := s.permissionRepo.Update(ctx, []any{
		clause.Eq{Column: "id", Value: permission.ID},
	}, map[string]any{
		"name": req.Name,
	}); err != nil {
		s.logger.Error("Failed to update permission", "error", err)
		return response.ResultErrDb(err)
	}

	permission.Name = req.Name
	return response.ResultSuccess("Permission updated successfully", permission)
}
