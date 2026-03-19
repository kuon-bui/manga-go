package gorm

import (
	"base-go/internal/pkg/config"
	"base-go/internal/pkg/logger"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
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
	if isProduction := cfg.Production; !isProduction {
		dsn += " sslmode=disable"
	}

	logger.Info("PostgreSQL dsn: ", dsn)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
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
