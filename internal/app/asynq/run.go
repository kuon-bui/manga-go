package asynq

import (
	"base-go/internal/pkg/common"
	"base-go/internal/pkg/config"
	"base-go/internal/pkg/logger"
	"base-go/internal/pkg/tracer"
	"context"

	"github.com/hibiken/asynq"
	"go.uber.org/fx"
)

type RunAsynqParams struct {
	fx.In
	Lc             fx.Lifecycle
	Server         *asynq.Server
	Mux            *asynq.ServeMux
	WorkerManagers []common.WorkerManager `group:"workerManagers"`
	Logger         *logger.Logger
	Config         *config.Config
}

func RunServer(p RunAsynqParams) {
	var wg []common.WorkerManager
	var cleanupTracer func(context.Context) error

	p.Lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			wg = p.WorkerManagers
			cleanupTracer = tracer.InitTracer(p.Config, p.Logger)
			p.Logger.Info("Starting asynq: ", len(wg))
			for _, w := range wg {
				w.RegisterWorkers()
			}

			go func() {
				if err := p.Server.Run(p.Mux); err != nil {
					p.Logger.Fatalf("could not run server: %v", err)
				}
			}()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			if cleanupTracer != nil {
				cleanupTracer(ctx)
			}

			p.Server.Shutdown()

			return nil
		},
	})
}
