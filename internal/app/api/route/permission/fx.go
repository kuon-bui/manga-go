package permissionroute

import (
	"manga-go/internal/app/api/common"

	"go.uber.org/fx"
)

var Module = fx.Module(
	"permission-route",
	common.ProvideAsRoute(NewPermissionRoute),
	fx.Provide(NewPermissionHandler),
)
