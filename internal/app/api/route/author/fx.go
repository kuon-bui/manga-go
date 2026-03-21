package authorroute

import (
	"manga-go/internal/app/api/common"

	"go.uber.org/fx"
)

var Module = fx.Module(
	"author-route",
	common.ProvideAsRoute(NewAuthorRoute),
	fx.Provide(NewAuthorHandler),
)
