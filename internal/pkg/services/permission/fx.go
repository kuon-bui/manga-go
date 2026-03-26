package permissionservice

import "go.uber.org/fx"

var Module = fx.Module(
	"permission-service",
	fx.Provide(NewPermissionService),
)
