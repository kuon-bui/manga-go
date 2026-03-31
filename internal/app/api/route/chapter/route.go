package chapterhandler

import (
	authmiddleware "manga-go/internal/app/middleware/auth"
	slugmiddleware "manga-go/internal/app/middleware/slug"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

type ChapterRoute struct {
	handler        *ChapterHandler
	r              *gin.Engine
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
		handler:        p.Handler,
		r:              p.R,
		authMiddleware: p.AuthMiddleware,
		slugMiddleware: p.SlugMiddleware,
	}
}

func (cr *ChapterRoute) Setup() {
	rg := cr.r.Group("/comics/:comicSlug/chapters", cr.authMiddleware.RequireJwt, cr.slugMiddleware.ResolveComicID)

	rg.GET("", cr.handler.listChapters)
	rg.GET("/:chapterSlug", cr.slugMiddleware.ResolveChapterID, cr.handler.getChapter)
	rg.POST("", cr.handler.createChapter)
	rg.PUT("/:chapterSlug", cr.slugMiddleware.ResolveChapterID, cr.handler.updateChapter)
	rg.PATCH("/:chapterSlug/publish", cr.slugMiddleware.ResolveChapterID, cr.handler.publishChapter)
}
