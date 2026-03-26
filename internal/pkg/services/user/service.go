package userservice

import (
	"manga-go/internal/pkg/config"
	jwtprovider "manga-go/internal/pkg/jwt_provider"
	"manga-go/internal/pkg/logger"
	rolerepo "manga-go/internal/pkg/repo/role"
	userrepo "manga-go/internal/pkg/repo/user"

	"github.com/hibiken/asynq"
	"go.uber.org/fx"
)

type UserService struct {
	logger      *logger.Logger
	userRepo    *userrepo.UserRepository
	jwtProvider *jwtprovider.JwtProvider
	config      *config.Config
	asynqClient *asynq.Client
	roleRepo    *rolerepo.RoleRepo
}

type UserServiceParams struct {
	fx.In

	Config      *config.Config
	Logger      *logger.Logger
	JwtProvider *jwtprovider.JwtProvider
	UserRepo    *userrepo.UserRepository
	AsynqClient *asynq.Client
	RoleRepo    *rolerepo.RoleRepo
}

func NewUserService(p UserServiceParams) *UserService {
	return &UserService{
		logger:      p.Logger,
		userRepo:    p.UserRepo,
		jwtProvider: p.JwtProvider,
		config:      p.Config,
		asynqClient: p.AsynqClient,
		roleRepo:    p.RoleRepo,
	}
}
