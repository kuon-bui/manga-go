package roleservice

import (
	"manga-go/internal/pkg/authorization"
	"manga-go/internal/pkg/logger"
	rolerepo "manga-go/internal/pkg/repo/role"

	"go.uber.org/fx"
)

type RoleService struct {
	logger        *logger.Logger
	roleRepo      *rolerepo.RoleRepo
	policyManager *authorization.PolicyManager
}

type RoleServiceParams struct {
	fx.In
	Logger        *logger.Logger
	RoleRepo      *rolerepo.RoleRepo
	PolicyManager *authorization.PolicyManager
}

func NewRoleService(params RoleServiceParams) *RoleService {
	return &RoleService{
		logger:        params.Logger,
		roleRepo:      params.RoleRepo,
		policyManager: params.PolicyManager,
	}
}
