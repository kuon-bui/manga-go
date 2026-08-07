package authorizationadmin

import (
	"manga-go/internal/pkg/authorization"
	"manga-go/internal/pkg/logger"
	authorizationrevision "manga-go/internal/pkg/repo/authorization_revision"
	rolerepo "manga-go/internal/pkg/repo/role"
	userrepo "manga-go/internal/pkg/repo/user"

	"go.uber.org/fx"
)

type ServiceParams struct {
	fx.In

	Logger        *logger.Logger
	RoleRepo      *rolerepo.RoleRepo
	UserRepo      *userrepo.UserRepository
	PolicyManager *authorization.PolicyManager
	Authorizer    *authorization.Authorizer
	Revisions     *authorizationrevision.Repo
	Cache         *RedisProfileCache
}

func NewService(params ServiceParams) *Service {
	return &Service{
		logger:        params.Logger,
		roleRepo:      params.RoleRepo,
		userRepo:      params.UserRepo,
		policyManager: params.PolicyManager,
		authorizer:    params.Authorizer,
		revisions:     params.Revisions,
		cache:         params.Cache,
	}
}

var Module = fx.Module(
	"authorization-admin-service",
	fx.Provide(NewRedisProfileCache, NewService),
)
