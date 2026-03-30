package casbin

import (
	"manga-go/internal/pkg/logger"

	"go.uber.org/fx"
)

var Module = fx.Module(
	"casbin",
	fx.Provide(NewEnforcer),
	fx.Invoke(func(enforcer *Enforcer, log *logger.Logger) {
		SeedGlobalPolicies(enforcer, log)
	}),
)
