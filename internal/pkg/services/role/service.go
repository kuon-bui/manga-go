package roleservice

import (
	"manga-go/internal/pkg/logger"
	permissionrepo "manga-go/internal/pkg/repo/permission"
	rolerepo "manga-go/internal/pkg/repo/role"

	"go.uber.org/fx"
)

type RoleService struct {
	logger         *logger.Logger
	roleRepo       *rolerepo.RoleRepo
	permissionRepo *permissionrepo.PermissionRepo
}

type RoleServiceParams struct {
	fx.In
	Logger         *logger.Logger
	RoleRepo       *rolerepo.RoleRepo
	PermissionRepo *permissionrepo.PermissionRepo
}

func NewRoleService(params RoleServiceParams) *RoleService {
	return &RoleService{
		logger:         params.Logger,
		roleRepo:       params.RoleRepo,
		permissionRepo: params.PermissionRepo,
	}
}
