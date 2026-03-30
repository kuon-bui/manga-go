package app

import (
	authmiddleware "manga-go/internal/app/middleware/auth"
	casbinmiddleware "manga-go/internal/app/middleware/casbin"
	slugmiddleware "manga-go/internal/app/middleware/slug"
	asynqclient "manga-go/internal/pkg/asynq"
	casbinpkg "manga-go/internal/pkg/casbin"
	"manga-go/internal/pkg/config"
	"manga-go/internal/pkg/gorm"
	jwtprovider "manga-go/internal/pkg/jwt_provider"
	"manga-go/internal/pkg/logger"
	"manga-go/internal/pkg/mail"
	objectstorage "manga-go/internal/pkg/object_storage"
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
		slugmiddleware.NewSlugMiddleware,
		casbinmiddleware.NewCasbinMiddleware,
		mail.NewMailDialer,
		asynqclient.NewAsynqClient,
		objectstorage.NewObjectStorage,
	),
	casbinpkg.Module,
	repo.Module,
	services.Module,
)
