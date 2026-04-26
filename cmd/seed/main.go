package main

import (
	"context"
	"manga-go/internal/pkg/config"
	gormdb "manga-go/internal/pkg/gorm"
	"manga-go/internal/pkg/logger"
	"manga-go/internal/pkg/redis"
	"manga-go/internal/pkg/repo"
	"manga-go/internal/pkg/seeder"
	"os"
	"time"

	"go.uber.org/fx"
)

func main() {
	log := logger.NewLogger()

	var runner *seeder.SeederRunner

	app := fx.New(
		fx.Provide(
			config.LoadConfig(true),
			logger.NewLogger,
			gormdb.ConnectGORM,
			redis.ConnectRedis,
			redis.NewRedis,
		),
		repo.Module,
		seeder.Module,
		fx.Populate(&runner),
		fx.NopLogger,
	)

	startCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := app.Start(startCtx); err != nil {
		log.Fatalf("Failed to start app: %v", err)
	}

	if err := runner.RunAll(context.Background()); err != nil {
		log.Errorf("Seeder failed: %v", err)

		stopCtx, stopCancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer stopCancel()
		_ = app.Stop(stopCtx)
		os.Exit(1)
	}

	stopCtx, stopCancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer stopCancel()
	if err := app.Stop(stopCtx); err != nil {
		log.Errorf("Failed to stop app: %v", err)
	}
}
