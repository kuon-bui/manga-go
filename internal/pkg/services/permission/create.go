package permissionservice

import (
	"context"
	"manga-go/internal/app/api/common/response"
	"manga-go/internal/pkg/model"
	permissionrequest "manga-go/internal/pkg/request/permission"
)

func (s *PermissionService) CreatePermission(ctx context.Context, req *permissionrequest.CreatePermissionRequest) response.Result {
	permission := model.Permission{
		Name: req.Name,
	}

	if err := s.permissionRepo.Create(ctx, &permission); err != nil {
		s.logger.Error("Failed to create permission", "error", err)
		return response.ResultErrDb(err)
	}

	return response.ResultSuccess("Permission created successfully", permission)
}
