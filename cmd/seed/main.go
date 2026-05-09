package main

import (
	"context"
	"flag"
	"manga-go/internal/pkg/config"
	gormdb "manga-go/internal/pkg/gorm"
	"manga-go/internal/pkg/logger"
	"manga-go/internal/pkg/redis"
	"manga-go/internal/pkg/repo"
	"manga-go/internal/pkg/seeder"
	"os"
	"time"

	"github.com/jaswdr/faker/v2"
	"go.uber.org/fx"
)

func main() {
	log := logger.NewLogger()
	truncateBeforeSeed := flag.Bool("truncate", false, "truncate seeded tables before running the seeder")
	flag.Parse()

	var runner *seeder.SeederRunner

	app := fx.New(
		fx.Provide(
			config.LoadConfig("config.seeder.yml"),
			logger.NewLogger,
			gormdb.ConnectGORM,
			redis.ConnectRedis,
			redis.NewRedis,
			faker.New,
		),
		repo.Module,
		seeder.Module,
		fx.Populate(&runner),
		fx.NopLogger,
	)

	startCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)

	if err := app.Start(startCtx); err != nil {
		cancel()
		log.Errorf("Failed to start app: %v", err)
		os.Exit(1)
	}
	cancel()
	if *truncateBeforeSeed {
		if err := runner.TruncateAll(context.Background()); err != nil {
			log.Errorf("Seeder truncate failed: %v", err)

			stopCtx, stopCancel := context.WithTimeout(context.Background(), 15*time.Second)
			_ = app.Stop(stopCtx)
			stopCancel()
			os.Exit(1)
		}
	}

	if err := runner.RunAll(context.Background()); err != nil {
		log.Errorf("Seeder failed: %v", err)

		stopCtx, stopCancel := context.WithTimeout(context.Background(), 15*time.Second)
		_ = app.Stop(stopCtx)
		stopCancel()
		os.Exit(1)
	}

	stopCtx, stopCancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer stopCancel()
	if err := app.Stop(stopCtx); err != nil {
		log.Errorf("Failed to stop app: %v", err)
	}
}
