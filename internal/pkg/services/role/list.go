package roleservice

import (
	"context"
	"manga-go/internal/app/api/common/response"
	"manga-go/internal/pkg/common"
)

func (s *RoleService) ListRoles(ctx context.Context, paging *common.Paging) response.Result {
	roles, total, err := s.roleRepo.FindPaginated(ctx, nil, paging, nil)
	if err != nil {
		s.logger.Error("Failed to list roles", "error", err)
		return response.ResultErrDb(err)
	}

	return response.ResultSuccess("Roles retrieved successfully", response.ResponsePaginationData(roles, total))
}
