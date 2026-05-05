package comicroute

import (
	authmiddleware "manga-go/internal/app/middleware/auth"
	slugmiddleware "manga-go/internal/app/middleware/slug"
	"manga-go/internal/pkg/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

type ComicRoute struct {
	*gin.Engine
	logger         *logger.Logger
	authMiddleware *authmiddleware.AuthMiddleware
	comicHandler   *ComicHandler
	slugMiddleware *slugmiddleware.SlugMiddleware
}

type ComicRouteParams struct {
	fx.In

	R              *gin.Engine
	Logger         *logger.Logger
	ComicHandler   *ComicHandler
	AuthMiddleware *authmiddleware.AuthMiddleware
	SlugMiddleware *slugmiddleware.SlugMiddleware
}

func NewComicRoute(params ComicRouteParams) *ComicRoute {
	return &ComicRoute{
		Engine:         params.R,
		logger:         params.Logger,
		authMiddleware: params.AuthMiddleware,
		comicHandler:   params.ComicHandler,
		slugMiddleware: params.SlugMiddleware,
	}
}

func (cr *ComicRoute) Setup() {
	rg := cr.Group("/comics", cr.authMiddleware.RequireJwt)

	rg.GET("", cr.comicHandler.getComics)
	rg.GET("/trending", cr.comicHandler.getTrendingComics)
	rg.POST("", cr.comicHandler.createComic)

	slugRg := rg.Group("/:comicSlug", cr.slugMiddleware.ResolveComicID)
	slugRg.GET("", cr.comicHandler.getComic)
	slugRg.POST("/follow", cr.comicHandler.followComic)
	slugRg.GET("/follow-status", cr.comicHandler.getComicFollowStatus)
	slugRg.PATCH("/follow-status", cr.comicHandler.updateComicFollowStatus)
	slugRg.PUT("", cr.comicHandler.updateComic)
	slugRg.PATCH("/status", cr.comicHandler.updateComicStatus)
	slugRg.PATCH("/publish", cr.comicHandler.publishComic)
	slugRg.DELETE("/follow", cr.comicHandler.unfollowComic)
	slugRg.DELETE("", cr.comicHandler.deleteComic)
}
