package roleservice

import (
	"context"
	"manga-go/internal/app/api/common/response"
)

func (s *RoleService) ListAllRoles(ctx context.Context) response.Result {
	roles, err := s.roleRepo.FindAll(ctx, nil, nil)
	if err != nil {
		s.logger.Error("Failed to list all roles", "error", err)
		return response.ResultErrDb(err)
	}

	return response.ResultSuccess("Roles retrieved successfully", roles)
}
