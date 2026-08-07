package authorizationadmin

import (
	"manga-go/internal/pkg/authorization"
	"manga-go/internal/pkg/logger"
	authorizationaudit "manga-go/internal/pkg/repo/authorization_audit"
	authorizationrevision "manga-go/internal/pkg/repo/authorization_revision"
	rolerepo "manga-go/internal/pkg/repo/role"
	userrepo "manga-go/internal/pkg/repo/user"

	"go.uber.org/fx"
	"gorm.io/gorm"
)

type ServiceParams struct {
	fx.In

	Logger        *logger.Logger
	DB            *gorm.DB
	RoleRepo      *rolerepo.RoleRepo
	UserRepo      *userrepo.UserRepository
	PolicyManager *authorization.PolicyManager
	Authorizer    *authorization.Authorizer
	AuditRepo     *authorizationaudit.Repo
	Revisions     *authorizationrevision.Repo
	Cache         *RedisProfileCache
}

func NewService(params ServiceParams) *Service {
	return &Service{
		logger:        params.Logger,
		db:            params.DB,
		roleRepo:      params.RoleRepo,
		userRepo:      params.UserRepo,
		policyManager: params.PolicyManager,
		authorizer:    params.Authorizer,
		auditRepo:     params.AuditRepo,
		revisions:     params.Revisions,
		cache:         params.Cache,
		locker:        NewPostgresMutationLocker(params.DB),
	}
}

var Module = fx.Module(
	"authorization-admin-service",
	fx.Provide(NewRedisProfileCache, NewService),
)
