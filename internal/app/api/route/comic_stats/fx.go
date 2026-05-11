package comicstatsroute

import (
	"manga-go/internal/app/api/common"

	"go.uber.org/fx"
)

var Module = fx.Module(
	"comic-stats-route",
	common.ProvideAsRoute(NewComicStatsRoute),
	fx.Provide(NewComicStatsHandler),
)
