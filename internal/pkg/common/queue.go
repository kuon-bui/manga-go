package common

import "go.uber.org/fx"

type WorkerManager interface {
	// RegisterWorkers registers all workers for this module
	RegisterWorkers()
}

func AsWorkerManager(f any) fx.Option {
	return fx.Provide(
		fx.Annotate(
			f,
			fx.As(new(WorkerManager)),
			fx.ResultTags(`group:"workerManagers"`),
		),
	)
}
