package app

import (
	authmiddleware "base-go/internal/app/middleware/auth"
	asynqclient "base-go/internal/pkg/asynq"
	"base-go/internal/pkg/config"
	"base-go/internal/pkg/gorm"
	jwtprovider "base-go/internal/pkg/jwt_provider"
	"base-go/internal/pkg/logger"
	"base-go/internal/pkg/mail"
	"base-go/internal/pkg/redis"
	"base-go/internal/pkg/repo"
	"base-go/internal/pkg/services"
	"base-go/internal/pkg/tracer"

	"go.uber.org/fx"
)

var Module = fx.Module(
	"app",
	fx.Provide(
		config.LoadConfig,
		gorm.ConnectGORM,
		redis.ConnectRedis,
		redis.NewRedis,
		tracer.InitTracer,
		logger.NewLogger,
		jwtprovider.NewJwtProvider,
		authmiddleware.NewAuthMiddleware,
		mail.NewMailDialer,
		asynqclient.NewAsynqClient,
	),
	repo.Module,
	services.Module,
)
