package comicstatsroute

import (
	comicstatsservice "manga-go/internal/pkg/services/comic_stats"

	"go.uber.org/fx"
)

type ComicStatsHandler struct {
	comicStatsService *comicstatsservice.ComicStatsService
}

type ComicStatsHandlerParams struct {
	fx.In
	ComicStatsService *comicstatsservice.ComicStatsService
}

func NewComicStatsHandler(p ComicStatsHandlerParams) *ComicStatsHandler {
	return &ComicStatsHandler{
		comicStatsService: p.ComicStatsService,
	}
}
