package permissionrepo

import "go.uber.org/fx"

var Module = fx.Module(
	"permission-repo",
	fx.Provide(NewPermissionRepo),
)
