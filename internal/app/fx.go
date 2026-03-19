package app

import (
	authmiddleware "manga-go/internal/app/middleware/auth"
	asynqclient "manga-go/internal/pkg/asynq"
	"manga-go/internal/pkg/config"
	"manga-go/internal/pkg/gorm"
	jwtprovider "manga-go/internal/pkg/jwt_provider"
	"manga-go/internal/pkg/logger"
	"manga-go/internal/pkg/mail"
	"manga-go/internal/pkg/redis"
	"manga-go/internal/pkg/repo"
	"manga-go/internal/pkg/services"
	"manga-go/internal/pkg/tracer"

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
