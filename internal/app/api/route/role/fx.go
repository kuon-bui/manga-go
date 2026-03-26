package roleroute

import (
	"manga-go/internal/app/api/common"

	"go.uber.org/fx"
)

var Module = fx.Module(
	"role-route",
	common.ProvideAsRoute(NewRoleRoute),
	fx.Provide(NewRoleHandler),
)
