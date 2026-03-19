package tracer

import (
	"manga-go/internal/pkg/config"

	"go.uber.org/fx"
)

var Module = fx.Module(
	"tracer",
	fx.Provide(
		NewTracer,
	),
	fx.Invoke(InitTracer),
)

func NewTracer(cfg *config.Config) (*Tracer, error) {
	return &Tracer{}, nil
}
