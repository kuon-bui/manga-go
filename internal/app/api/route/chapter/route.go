package chapterhandler

import (
	authmiddleware "manga-go/internal/app/middleware/auth"
	slugmiddleware "manga-go/internal/app/middleware/slug"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

type ChapterRoute struct {
	*gin.Engine
	handler        *ChapterHandler
	authMiddleware *authmiddleware.AuthMiddleware
	slugMiddleware *slugmiddleware.SlugMiddleware
}

type ChapterRouteParams struct {
	fx.In

	Handler        *ChapterHandler
	R              *gin.Engine
	AuthMiddleware *authmiddleware.AuthMiddleware
	SlugMiddleware *slugmiddleware.SlugMiddleware
}

func NewChapterRoute(p ChapterRouteParams) *ChapterRoute {
	return &ChapterRoute{
		Engine:         p.R,
		handler:        p.Handler,
		authMiddleware: p.AuthMiddleware,
		slugMiddleware: p.SlugMiddleware,
	}
}

func (cr *ChapterRoute) Setup() {
	cr.GET("/chapters/recent-updates", cr.authMiddleware.RequireJwt, cr.handler.getRecentUpdates)

	rg := cr.Group("/comics/:comicSlug/chapters", cr.authMiddleware.RequireJwt, cr.slugMiddleware.ResolveComicID)

	rg.GET("", cr.handler.listChapters)
	rg.POST("", cr.handler.createChapter)

	rgSlug := rg.Group("/:chapterSlug", cr.slugMiddleware.ResolveChapterID)

	rgSlug.GET("", cr.handler.getChapter)
	rgSlug.PUT("", cr.handler.updateChapter)
	rgSlug.PUT("/pages", cr.handler.updateChapterPages)
	rgSlug.PATCH("/publish", cr.handler.publishChapter)
	rgSlug.PATCH("/mark-as-read", cr.handler.markChapterAsRead)

	readingProgressRg := rgSlug.Group("/reading-progress")
	readingProgressRg.GET("", cr.handler.getReadingProgress)
	readingProgressRg.PATCH("", cr.handler.updateReadingProgress)
}
