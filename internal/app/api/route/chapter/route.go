package chapterhandler

import (
	authmiddleware "manga-go/internal/app/middleware/auth"
	authzmiddleware "manga-go/internal/app/middleware/authz"
	slugmiddleware "manga-go/internal/app/middleware/slug"
	"manga-go/internal/pkg/authorization"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

type ChapterRoute struct {
	*gin.Engine
	handler         *ChapterHandler
	authMiddleware  *authmiddleware.AuthMiddleware
	authzMiddleware *authzmiddleware.AuthzMiddleware
	slugMiddleware  *slugmiddleware.SlugMiddleware
}

type ChapterRouteParams struct {
	fx.In

	Handler         *ChapterHandler
	R               *gin.Engine
	AuthMiddleware  *authmiddleware.AuthMiddleware
	AuthzMiddleware *authzmiddleware.AuthzMiddleware
	SlugMiddleware  *slugmiddleware.SlugMiddleware
}

func NewChapterRoute(p ChapterRouteParams) *ChapterRoute {
	return &ChapterRoute{
		Engine:          p.R,
		handler:         p.Handler,
		authMiddleware:  p.AuthMiddleware,
		authzMiddleware: p.AuthzMiddleware,
		slugMiddleware:  p.SlugMiddleware,
	}
}

func (cr *ChapterRoute) Setup() {
	requirePublishedChapterRead := authzmiddleware.Require(cr.authzMiddleware, authorization.ActionRead, authorization.ObjectChapter, authzmiddleware.Published())
	cr.GET("/chapters/recent-updates", cr.authMiddleware.OptionalJwt, requirePublishedChapterRead, cr.handler.getRecentUpdates)

	publicRg := cr.Group("/comics/:comicSlug/chapters", cr.authMiddleware.OptionalJwt, cr.slugMiddleware.ResolveComicID)
	requireComicRead := authzmiddleware.Require(cr.authzMiddleware, authorization.ActionRead, authorization.ObjectComic, cr.authzMiddleware.Comic())
	requireChapterRead := authzmiddleware.Require(cr.authzMiddleware, authorization.ActionRead, authorization.ObjectChapter, cr.authzMiddleware.Chapter(), cr.authzMiddleware.ComicGroupFromContext())
	publicRg.GET("", requireComicRead, requirePublishedChapterRead, cr.handler.listChapters)
	publicRg.GET("/:chapterSlug", requireComicRead, cr.slugMiddleware.ResolveChapterID, requireChapterRead, cr.handler.getChapter)

	rg := cr.Group("/comics/:comicSlug/chapters", cr.authMiddleware.RequireJwt, cr.slugMiddleware.ResolveComicID)
	requireChapterCreate := authzmiddleware.Require(cr.authzMiddleware, authorization.ActionCreate, authorization.ObjectChapter, cr.authzMiddleware.ComicGroupFromContext())
	requireChapterUpdate := authzmiddleware.Require(cr.authzMiddleware, authorization.ActionUpdate, authorization.ObjectChapter, cr.authzMiddleware.Chapter(), cr.authzMiddleware.ComicGroupFromContext())
	requireChapterPublish := authzmiddleware.Require(cr.authzMiddleware, authorization.ActionPublish, authorization.ObjectChapter, cr.authzMiddleware.Chapter(), cr.authzMiddleware.ComicGroupFromContext())

	rg.POST("", requireChapterCreate, cr.handler.createChapter)

	rgSlug := rg.Group("/:chapterSlug", cr.slugMiddleware.ResolveChapterID)

	rgSlug.PUT("", requireChapterUpdate, cr.handler.updateChapter)
	rgSlug.PUT("/pages", requireChapterUpdate, cr.handler.updateChapterPages)
	rgSlug.PATCH("/publish", requireChapterPublish, cr.handler.publishChapter)
	rgSlug.PATCH("/mark-as-read", cr.handler.markChapterAsRead)

	readingProgressRg := rgSlug.Group("/reading-progress")
	readingProgressRg.GET("", cr.handler.getReadingProgress)
	readingProgressRg.PATCH("", cr.handler.updateReadingProgress)
}
