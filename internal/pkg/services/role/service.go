package roleservice

import (
	"manga-go/internal/pkg/authorization"
	"manga-go/internal/pkg/logger"
	rolerepo "manga-go/internal/pkg/repo/role"
	authorizationadmin "manga-go/internal/pkg/services/authorization_admin"

	"go.uber.org/fx"
)

type RoleService struct {
	logger        *logger.Logger
	roleRepo      *rolerepo.RoleRepo
	policyManager *authorization.PolicyManager
	authAdmin     *authorizationadmin.Service
}

type RoleServiceParams struct {
	fx.In
	Logger        *logger.Logger
	RoleRepo      *rolerepo.RoleRepo
	PolicyManager *authorization.PolicyManager
	AuthAdmin     *authorizationadmin.Service
}

func NewRoleService(params RoleServiceParams) *RoleService {
	return &RoleService{
		logger:        params.Logger,
		roleRepo:      params.RoleRepo,
		policyManager: params.PolicyManager,
		authAdmin:     params.AuthAdmin,
	}
}
