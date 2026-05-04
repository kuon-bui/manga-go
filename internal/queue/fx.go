package queue

import (
	asynqclient "manga-go/internal/pkg/asynq"
	"manga-go/internal/pkg/config"
	"manga-go/internal/pkg/gorm"
	"manga-go/internal/pkg/logger"
	"manga-go/internal/pkg/mail"
	objectstorage "manga-go/internal/pkg/object_storage"
	"manga-go/internal/pkg/redis"
	"manga-go/internal/pkg/repo"
	fileservice "manga-go/internal/pkg/services/file"
	notificationservice "manga-go/internal/pkg/services/notification"
	"manga-go/internal/pkg/tracer"
	"manga-go/internal/queue/asynq"

	"go.uber.org/fx"
)

var Module = fx.Module(
	"queue",
	fx.Provide(
		config.LoadConfig("config.yml"),
		gorm.ConnectGORM,
		redis.ConnectRedis,
		redis.NewRedis,
		tracer.InitTracer,
		logger.NewLogger,
		mail.NewMailDialer,
		asynqclient.NewAsynqClient,
		objectstorage.NewObjectStorage,
	),
	repo.Module,
	fileservice.Module,
	notificationservice.Module,
	asynq.Module,
)
