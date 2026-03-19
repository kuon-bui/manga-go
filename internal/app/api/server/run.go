package server

import (
	"base-go/internal/app/api/common"
	"base-go/internal/pkg/config"
	"base-go/internal/pkg/logger"
	"base-go/internal/pkg/tracer"
	"context"
	"net/http"

	"go.uber.org/fx"
	"gorm.io/gorm"
)

type RunServerParams struct {
	fx.In

	Lc     fx.Lifecycle
	Config *config.Config
	Logger *logger.Logger
	Gorm   *gorm.DB
	Server *http.Server
	Routes []common.Route `group:"routes"`
}

func RunServer(p RunServerParams) {
	var cleanupTracer func(context.Context) error

	p.Lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			cleanupTracer = tracer.InitTracer(p.Config, p.Logger)

			p.Logger.Info("Starting server...")
			for _, route := range p.Routes {
				route.Setup()
			}
			return nil
		},
		OnStop: func(ctx context.Context) error {
			p.Logger.Info("Stopping server...")

			// Graceful shutdown
			if err := p.Server.Shutdown(ctx); err != nil {
				p.Logger.Fatal("Server forced to shutdown: ", err)
			}

			// Waiting for the goroutines to have a chance to complete
			// time.Sleep(3 * time.Second)
			if cleanupTracer != nil {
				cleanupTracer(ctx)
			}

			return nil
		},
	})
}
