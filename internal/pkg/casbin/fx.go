package casbin

import "go.uber.org/fx"

var Module = fx.Module(
	"casbin",
	fx.Provide(NewEnforcer),
)
