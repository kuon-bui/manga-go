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
	cr.GET("/chapters/recent-updates", cr.authMiddleware.RequireJwt, cr.handler.getRecentUpdates)

	rg := cr.Group("/comics/:comicSlug/chapters", cr.authMiddleware.RequireJwt, cr.slugMiddleware.ResolveComicID)
	requireChapterCreate := authzmiddleware.Require(cr.authzMiddleware, authorization.ActionCreate, authorization.ObjectChapter, cr.authzMiddleware.ComicGroupFromContext())
	requireChapterRead := authzmiddleware.Require(cr.authzMiddleware, authorization.ActionRead, authorization.ObjectChapter, cr.authzMiddleware.Chapter(), cr.authzMiddleware.ComicGroupFromContext())
	requireChapterUpdate := authzmiddleware.Require(cr.authzMiddleware, authorization.ActionUpdate, authorization.ObjectChapter, cr.authzMiddleware.Chapter(), cr.authzMiddleware.ComicGroupFromContext())
	requireChapterPublish := authzmiddleware.Require(cr.authzMiddleware, authorization.ActionPublish, authorization.ObjectChapter, cr.authzMiddleware.Chapter(), cr.authzMiddleware.ComicGroupFromContext())

	rg.GET("", cr.handler.listChapters)
	rg.POST("", requireChapterCreate, cr.handler.createChapter)

	rgSlug := rg.Group("/:chapterSlug", cr.slugMiddleware.ResolveChapterID)

	rgSlug.GET("", requireChapterRead, cr.handler.getChapter)
	rgSlug.PUT("", requireChapterUpdate, cr.handler.updateChapter)
	rgSlug.PUT("/pages", requireChapterUpdate, cr.handler.updateChapterPages)
	rgSlug.PATCH("/publish", requireChapterPublish, cr.handler.publishChapter)
	rgSlug.PATCH("/mark-as-read", cr.handler.markChapterAsRead)

	readingProgressRg := rgSlug.Group("/reading-progress")
	readingProgressRg.GET("", cr.handler.getReadingProgress)
	readingProgressRg.PATCH("", cr.handler.updateReadingProgress)
}
