package comicroute

import (
	authmiddleware "manga-go/internal/app/middleware/auth"
	authzmiddleware "manga-go/internal/app/middleware/authz"
	slugmiddleware "manga-go/internal/app/middleware/slug"
	"manga-go/internal/pkg/authorization"
	"manga-go/internal/pkg/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

type ComicRoute struct {
	*gin.Engine
	logger          *logger.Logger
	authMiddleware  *authmiddleware.AuthMiddleware
	authzMiddleware *authzmiddleware.AuthzMiddleware
	comicHandler    *ComicHandler
	slugMiddleware  *slugmiddleware.SlugMiddleware
}

type ComicRouteParams struct {
	fx.In

	R               *gin.Engine
	Logger          *logger.Logger
	ComicHandler    *ComicHandler
	AuthMiddleware  *authmiddleware.AuthMiddleware
	AuthzMiddleware *authzmiddleware.AuthzMiddleware
	SlugMiddleware  *slugmiddleware.SlugMiddleware
}

func NewComicRoute(params ComicRouteParams) *ComicRoute {
	return &ComicRoute{
		Engine:          params.R,
		logger:          params.Logger,
		authMiddleware:  params.AuthMiddleware,
		authzMiddleware: params.AuthzMiddleware,
		comicHandler:    params.ComicHandler,
		slugMiddleware:  params.SlugMiddleware,
	}
}

func (cr *ComicRoute) Setup() {
	publicRg := cr.Group("/comics", cr.authMiddleware.OptionalJwt)
	requireComicRead := authzmiddleware.Require(cr.authzMiddleware, authorization.ActionRead, authorization.ObjectComic, cr.authzMiddleware.Comic())
	requirePublishedComicRead := authzmiddleware.Require(cr.authzMiddleware, authorization.ActionRead, authorization.ObjectComic, authzmiddleware.Published())
	requireComicCreate := authzmiddleware.Require(cr.authzMiddleware, authorization.ActionCreate, authorization.ObjectComic, cr.authzMiddleware.CurrentUserTranslationGroup())
	requireComicUpdate := authzmiddleware.Require(cr.authzMiddleware, authorization.ActionUpdate, authorization.ObjectComic, cr.authzMiddleware.Comic())
	requireComicDelete := authzmiddleware.Require(cr.authzMiddleware, authorization.ActionDelete, authorization.ObjectComic, cr.authzMiddleware.Comic())
	requireComicPublish := authzmiddleware.Require(cr.authzMiddleware, authorization.ActionPublish, authorization.ObjectComic, cr.authzMiddleware.Comic())

	publicRg.GET("", requirePublishedComicRead, cr.comicHandler.getComics)
	publicRg.GET("/trending", requirePublishedComicRead, cr.comicHandler.getTrendingComics)
	publicRg.GET("/:comicSlug", cr.slugMiddleware.ResolveComicID, requireComicRead, cr.comicHandler.getComic)

	privateRg := cr.Group("/comics", cr.authMiddleware.RequireJwt)
	privateRg.POST("", requireComicCreate, cr.comicHandler.createComic)

	slugRg := privateRg.Group("/:comicSlug", cr.slugMiddleware.ResolveComicID)
	slugRg.POST("/follow", cr.comicHandler.followComic)
	slugRg.GET("/follow-status", cr.comicHandler.getComicFollowStatus)
	slugRg.PATCH("/follow-status", cr.comicHandler.updateComicFollowStatus)
	slugRg.PUT("", requireComicUpdate, cr.comicHandler.updateComic)
	slugRg.PATCH("/status", requireComicUpdate, cr.comicHandler.updateComicStatus)
	slugRg.PATCH("/publish", requireComicPublish, cr.comicHandler.publishComic)
	slugRg.DELETE("/follow", cr.comicHandler.unfollowComic)
	slugRg.DELETE("", requireComicDelete, cr.comicHandler.deleteComic)
}
