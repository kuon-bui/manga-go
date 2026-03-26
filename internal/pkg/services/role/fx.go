package roleservice

import "go.uber.org/fx"

var Module = fx.Module(
	"role-service",
	fx.Provide(NewRoleService),
)
