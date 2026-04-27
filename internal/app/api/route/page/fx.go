package pageroute

import (
	"manga-go/internal/app/api/common"

	"go.uber.org/fx"
)

var Module = fx.Module(
	"page-route",
	common.ProvideAsRoute(NewPageRoute),
	fx.Provide(NewPageHandler),
)
