package comicstatsroute

import (
	authmiddleware "manga-go/internal/app/middleware/auth"
	authzmiddleware "manga-go/internal/app/middleware/authz"
	"manga-go/internal/pkg/authorization"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

type ComicStatsRoute struct {
	*gin.Engine
	authMiddleware  *authmiddleware.AuthMiddleware
	authzMiddleware *authzmiddleware.AuthzMiddleware
	handler         *ComicStatsHandler
}

type ComicStatsRouteParams struct {
	fx.In
	R               *gin.Engine
	AuthMiddleware  *authmiddleware.AuthMiddleware
	AuthzMiddleware *authzmiddleware.AuthzMiddleware
	Handler         *ComicStatsHandler
}

func NewComicStatsRoute(params ComicStatsRouteParams) *ComicStatsRoute {
	return &ComicStatsRoute{
		Engine:          params.R,
		authMiddleware:  params.AuthMiddleware,
		authzMiddleware: params.AuthzMiddleware,
		handler:         params.Handler,
	}
}

func (r *ComicStatsRoute) Setup() {
	// Recomputing stats mutates every comic in scope, so it takes the same
	// platform-wide comic update right an editor needs — not merely a session.
	requireComicStatsTrigger := authzmiddleware.Require(
		r.authzMiddleware,
		authorization.ActionUpdate,
		authorization.ObjectComic,
	)

	rg := r.Group("/admin/comic-stats", r.authMiddleware.RequireJwt, requireComicStatsTrigger)

	rg.POST("/trigger/:id", r.handler.triggerComicStats)
	rg.POST("/trigger-all", r.handler.triggerAllComicStats)
}
