package roleservice

import (
	"context"
	"manga-go/internal/app/api/common/response"
	"manga-go/internal/pkg/model"
	rolerequest "manga-go/internal/pkg/request/role"
)

func (s *RoleService) CreateRole(ctx context.Context, req *rolerequest.CreateRoleRequest) response.Result {
	if s.authAdmin != nil && s.authAdmin.MutationReady() {
		return s.authAdmin.CreateRole(ctx, req.Name, req.Description)
	}
	role := model.Role{
		Name:        req.Name,
		Description: req.Description,
	}

	if err := s.roleRepo.Create(ctx, &role); err != nil {
		s.logger.Error("Failed to create role", "error", err)
		return response.ResultErrDb(err)
	}

	return response.ResultSuccess("Role created successfully", role)
}
