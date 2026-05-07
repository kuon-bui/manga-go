package authorization

import "go.uber.org/fx"

var Module = fx.Module(
	"authorization",
	fx.Provide(
		NewAuthorizer,
		NewPolicyManager,
	),
)
