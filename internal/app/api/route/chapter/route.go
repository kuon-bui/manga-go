package chapterhandler

import (
	authmiddleware "manga-go/internal/app/middleware/auth"
	comicmiddleware "manga-go/internal/app/middleware/comic"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

type ChapterRoute struct {
	handler         *ChapterHandler
	r               *gin.Engine
	authMiddleware  *authmiddleware.AuthMiddleware
	comicMiddleware *comicmiddleware.ComicMiddleware
}

type ChapterRouteParams struct {
	fx.In

	Handler         *ChapterHandler
	R               *gin.Engine
	AuthMiddleware  *authmiddleware.AuthMiddleware
	ComicMiddleware *comicmiddleware.ComicMiddleware
}

func NewChapterRoute(p ChapterRouteParams) *ChapterRoute {
	return &ChapterRoute{
		handler:         p.Handler,
		r:               p.R,
		authMiddleware:  p.AuthMiddleware,
		comicMiddleware: p.ComicMiddleware,
	}
}

func (cr *ChapterRoute) Setup() {
	rg := cr.r.Group("/comics/:comicSlug/chapters", cr.authMiddleware.RequireJwt, cr.comicMiddleware.ResolveComicID)

	rg.GET("", cr.handler.listChapters)
	rg.GET("/:chapterSlug", cr.handler.getChapter)
	rg.POST("", cr.handler.createChapter)
	rg.PUT("/:chapterSlug", cr.handler.updateChapter)
}
