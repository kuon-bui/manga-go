package comicstatsservice

import "go.uber.org/fx"

var Module = fx.Module(
	"comic-stats-service",
	fx.Provide(NewComicStatsService),
)
