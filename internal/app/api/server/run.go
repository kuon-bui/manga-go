package server

import (
	"context"
	"errors"
	"manga-go/internal/app/api/common"
	"manga-go/internal/pkg/config"
	"manga-go/internal/pkg/logger"
	"manga-go/internal/pkg/tracer"
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
			go func() {
				p.Logger.Infof("Starting HTTP server on port %d", p.Config.Service.Port)
				if err := p.Server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
					p.Logger.Errorf("HTTP server failed: %v", err)
				}
				p.Logger.Info("Server closed")
			}()
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
				if err := cleanupTracer(ctx); err != nil {
					p.Logger.Errorf("Failed to cleanup tracer: %v", err)
				}
			}

			return nil
		},
	})
}
