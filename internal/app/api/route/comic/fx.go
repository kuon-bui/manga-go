package comicroute

import (
	"manga-go/internal/app/api/common"

	"go.uber.org/fx"
)

var Module = fx.Module(
	"comic-route",
	common.ProvideAsRoute(NewComicRoute),
	fx.Provide(NewComicHandler),
)
