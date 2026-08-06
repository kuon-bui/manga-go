package permissionservice

import (
	"manga-go/internal/pkg/authorization"
	"manga-go/internal/pkg/logger"
	permissionrepo "manga-go/internal/pkg/repo/permission"

	"go.uber.org/fx"
)

type PermissionService struct {
	logger         *logger.Logger
	permissionRepo *permissionrepo.PermissionRepo
	policyManager  *authorization.PolicyManager
}

type PermissionServiceParams struct {
	fx.In
	Logger         *logger.Logger
	PermissionRepo *permissionrepo.PermissionRepo
	PolicyManager  *authorization.PolicyManager
}

func NewPermissionService(params PermissionServiceParams) *PermissionService {
	return &PermissionService{
		logger:         params.Logger,
		permissionRepo: params.PermissionRepo,
		policyManager:  params.PolicyManager,
	}
}
