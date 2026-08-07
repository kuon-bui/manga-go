package userservice

import (
	"manga-go/internal/pkg/authorization"
	"manga-go/internal/pkg/config"
	jwtprovider "manga-go/internal/pkg/jwt_provider"
	"manga-go/internal/pkg/logger"
	rolerepo "manga-go/internal/pkg/repo/role"
	userrepo "manga-go/internal/pkg/repo/user"
	authorizationadmin "manga-go/internal/pkg/services/authorization_admin"

	"github.com/hibiken/asynq"
	"go.uber.org/fx"
)

type UserService struct {
	logger        *logger.Logger
	userRepo      *userrepo.UserRepository
	jwtProvider   *jwtprovider.JwtProvider
	config        *config.Config
	asynqClient   *asynq.Client
	roleRepo      *rolerepo.RoleRepo
	policyManager *authorization.PolicyManager
	authAdmin     *authorizationadmin.Service
}

type UserServiceParams struct {
	fx.In

	Config        *config.Config
	Logger        *logger.Logger
	JwtProvider   *jwtprovider.JwtProvider
	UserRepo      *userrepo.UserRepository
	AsynqClient   *asynq.Client
	RoleRepo      *rolerepo.RoleRepo
	PolicyManager *authorization.PolicyManager
	AuthAdmin     *authorizationadmin.Service
}

func NewUserService(p UserServiceParams) *UserService {
	return &UserService{
		logger:        p.Logger,
		userRepo:      p.UserRepo,
		jwtProvider:   p.JwtProvider,
		config:        p.Config,
		asynqClient:   p.AsynqClient,
		roleRepo:      p.RoleRepo,
		policyManager: p.PolicyManager,
		authAdmin:     p.AuthAdmin,
	}
}
