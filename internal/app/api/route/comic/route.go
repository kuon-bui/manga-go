package comicroute

import (
	authmiddleware "manga-go/internal/app/middleware/auth"
	casbinmiddleware "manga-go/internal/app/middleware/casbin"
	"manga-go/internal/pkg/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

type ComicRoute struct {
	logger          *logger.Logger
	r               *gin.Engine
	authMiddleware  *authmiddleware.AuthMiddleware
	casbinMiddleware *casbinmiddleware.CasbinMiddleware
	comicHandler    *ComicHandler
}

type ComicRouteParams struct {
	fx.In

	R               *gin.Engine
	Logger          *logger.Logger
	ComicHandler    *ComicHandler
	AuthMiddleware  *authmiddleware.AuthMiddleware
	CasbinMiddleware *casbinmiddleware.CasbinMiddleware
}

func NewComicRoute(params ComicRouteParams) *ComicRoute {
	return &ComicRoute{
		logger:          params.Logger,
		r:               params.R,
		authMiddleware:  params.AuthMiddleware,
		casbinMiddleware: params.CasbinMiddleware,
		comicHandler:    params.ComicHandler,
	}
}

func (cr *ComicRoute) Setup() {
	rg := cr.r.Group("/comics")

	// Public read endpoints – no authentication required
	rg.GET("", cr.comicHandler.getComics)
	rg.GET("/:comicSlug", cr.comicHandler.getComic)

	// Write endpoints require authentication. Admin access is enforced at the service layer
	// via the Casbin admin role bypass (g(sub, "admin", "global") in the model matcher).
	auth := rg.Group("", cr.authMiddleware.RequireJwt)
	auth.POST("", cr.comicHandler.createComic)
	auth.PUT("/:comicSlug", cr.comicHandler.updateComic)
	auth.DELETE("/:comicSlug", cr.comicHandler.deleteComic)
}
