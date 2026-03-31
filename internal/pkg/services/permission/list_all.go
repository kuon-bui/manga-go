package permissionservice

import (
	"context"
	"manga-go/internal/app/api/common/response"
)

func (s *PermissionService) ListAllPermissions(ctx context.Context) response.Result {
	permissions, err := s.permissionRepo.FindAll(ctx, nil, nil)
	if err != nil {
		s.logger.Error("Failed to list all permissions", "error", err)
		return response.ResultErrDb(err)
	}

	return response.ResultSuccess("Permissions retrieved successfully", permissions)
}
