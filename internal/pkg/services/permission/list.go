package permissionservice

import (
	"context"
	"manga-go/internal/app/api/common/response"
	"manga-go/internal/pkg/common"
)

func (s *PermissionService) ListPermissions(ctx context.Context, paging *common.Paging) response.Result {
	permissions, total, err := s.permissionRepo.FindPaginated(ctx, nil, paging, nil)
	if err != nil {
		s.logger.Error("Failed to list permissions", "error", err)
		return response.ResultErrDb(err)
	}

	return response.ResultSuccess("Permissions retrieved successfully", response.ResponsePaginationData(permissions, total))
}
