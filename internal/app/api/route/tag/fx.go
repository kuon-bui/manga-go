package tagroute

import (
	"manga-go/internal/app/api/common"

	"go.uber.org/fx"
)

var Module = fx.Module(
	"tag-route",
	common.ProvideAsRoute(NewTagRoute),
	fx.Provide(NewTagHandler),
)
