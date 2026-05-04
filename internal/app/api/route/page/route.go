package pageroute

import (
	authmiddleware "manga-go/internal/app/middleware/auth"
	"manga-go/internal/pkg/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

type PageRoute struct {
	logger         *logger.Logger
	r              *gin.Engine
	authMiddleware *authmiddleware.AuthMiddleware
	handler        *PageHandler
}

type PageRouteParams struct {
	fx.In

	Logger         *logger.Logger
	R              *gin.Engine
	AuthMiddleware *authmiddleware.AuthMiddleware
	Handler        *PageHandler
}

func NewPageRoute(p PageRouteParams) *PageRoute {
	return &PageRoute{
		logger:         p.Logger,
		r:              p.R,
		authMiddleware: p.AuthMiddleware,
		handler:        p.Handler,
	}
}

func (pr *PageRoute) Setup() {
	rg := pr.r.Group("/pages", pr.authMiddleware.RequireJwt)

	rg.POST("/:pageId/reactions", pr.handler.handleReaction)
}
