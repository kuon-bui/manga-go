package comicroute

import (
	authmiddleware "manga-go/internal/app/middleware/auth"
	"manga-go/internal/pkg/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

type ComicRoute struct {
	logger         *logger.Logger
	r              *gin.Engine
	authMiddleware *authmiddleware.AuthMiddleware
	comicHandler   *ComicHandler
}

type ComicRouteParams struct {
	fx.In

	R              *gin.Engine
	Logger         *logger.Logger
	ComicHandler   *ComicHandler
	AuthMiddleware *authmiddleware.AuthMiddleware
}

func NewComicRoute(params ComicRouteParams) *ComicRoute {
	return &ComicRoute{
		logger:         params.Logger,
		r:              params.R,
		authMiddleware: params.AuthMiddleware,
		comicHandler:   params.ComicHandler,
	}
}

func (cr *ComicRoute) Setup() {
	rg := cr.r.Group("/comics", cr.authMiddleware.RequireJwt)

	rg.GET("", cr.comicHandler.getComics)
	rg.POST("", cr.comicHandler.createComic)

	slugRg := rg.Group("/:comicSlug")
	slugRg.GET("", cr.comicHandler.getComic)
	slugRg.PUT("", cr.comicHandler.updateComic)
	slugRg.PATCH("/status", cr.comicHandler.updateComicStatus)
	slugRg.PATCH("/publish", cr.comicHandler.publishComic)
	slugRg.DELETE("", cr.comicHandler.deleteComic)
}
