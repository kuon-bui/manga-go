package gorm

import (
	"fmt"
	"log"
	"manga-go/internal/pkg/config"
	"manga-go/internal/pkg/logger"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
	"gorm.io/plugin/opentelemetry/tracing"
)

func ConnectGORM(cfg *config.Config, logger *logger.Logger) *gorm.DB {
	postgresqlConfig := cfg.PostgreSQL
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d",
		postgresqlConfig.Host,
		postgresqlConfig.User,
		postgresqlConfig.Password,
		postgresqlConfig.Database,
		postgresqlConfig.Port,
	)

	// nếu như không phải môi trường Production thì thêm sslmode=disable vào dsn
	if cfg.RunMode != config.RunModeProduction {
		dsn += " sslmode=disable"
	}

	logger.Info("PostgreSQL dsn: ", dsn)

	gormLogger := gormlogger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		gormlogger.Config{
			SlowThreshold:             200 * time.Millisecond,
			IgnoreRecordNotFoundError: cfg.RunMode == config.RunModeSeeder,
			Colorful:                  true,
		},
	)

	gormConfig := &gorm.Config{
		Logger: gormLogger,
	}

	db, err := gorm.Open(postgres.Open(dsn), gormConfig)
	if err != nil {
		logger.Error("Connect to PostgreSQL error: ", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		logger.Error("Get sql.DB error: ", err)
		return nil
	}

	if err := sqlDB.Ping(); err != nil {
		logger.Error("Ping PostgreSQL error: ", err)
		return nil
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)

	if err := db.Use(tracing.NewPlugin()); err != nil {
		logger.Error("Use tracing plugin error: ", err)
	}

	logger.Info("Connected to PostgreSQL")

	return db
}
