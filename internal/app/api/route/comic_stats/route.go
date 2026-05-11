package comicstatsroute

import (
	authmiddleware "manga-go/internal/app/middleware/auth"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

type ComicStatsRoute struct {
	*gin.Engine
	authMiddleware *authmiddleware.AuthMiddleware
	handler        *ComicStatsHandler
}

type ComicStatsRouteParams struct {
	fx.In
	R              *gin.Engine
	AuthMiddleware *authmiddleware.AuthMiddleware
	Handler        *ComicStatsHandler
}

func NewComicStatsRoute(params ComicStatsRouteParams) *ComicStatsRoute {
	return &ComicStatsRoute{
		Engine:         params.R,
		authMiddleware: params.AuthMiddleware,
		handler:        params.Handler,
	}
}

func (r *ComicStatsRoute) Setup() {
	rg := r.Group("/admin/comic-stats", r.authMiddleware.RequireJwt)

	rg.POST("/trigger/:id", r.handler.triggerComicStats)
	rg.POST("/trigger-all", r.handler.triggerAllComicStats)
}
