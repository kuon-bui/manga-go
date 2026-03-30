package chapterhandler

import (
	authmiddleware "manga-go/internal/app/middleware/auth"
	casbinmiddleware "manga-go/internal/app/middleware/casbin"
	slugmiddleware "manga-go/internal/app/middleware/slug"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

type ChapterRoute struct {
	handler          *ChapterHandler
	r                *gin.Engine
	authMiddleware   *authmiddleware.AuthMiddleware
	casbinMiddleware *casbinmiddleware.CasbinMiddleware
	slugMiddleware   *slugmiddleware.SlugMiddleware
}

type ChapterRouteParams struct {
	fx.In

	Handler          *ChapterHandler
	R                *gin.Engine
	AuthMiddleware   *authmiddleware.AuthMiddleware
	CasbinMiddleware *casbinmiddleware.CasbinMiddleware
	SlugMiddleware   *slugmiddleware.SlugMiddleware
}

func NewChapterRoute(p ChapterRouteParams) *ChapterRoute {
	return &ChapterRoute{
		handler:          p.Handler,
		r:                p.R,
		authMiddleware:   p.AuthMiddleware,
		casbinMiddleware: p.CasbinMiddleware,
		slugMiddleware:   p.SlugMiddleware,
	}
}

func (cr *ChapterRoute) Setup() {
	rg := cr.r.Group("/comics/:comicSlug/chapters", cr.slugMiddleware.ResolveComicID)

	// Public read endpoints
	rg.GET("", cr.handler.listChapters)
	rg.GET("/:chapterSlug", cr.slugMiddleware.ResolveChapterID, cr.handler.getChapter)

	// Write endpoints require authentication and chapter permission within the comic's group
	auth := rg.Group("", cr.authMiddleware.RequireJwt)
	auth.POST("",
		cr.casbinMiddleware.RequireGroupPermission("chapter", "create"),
		cr.handler.createChapter,
	)
	auth.PUT("/:chapterSlug",
		cr.slugMiddleware.ResolveChapterID,
		cr.casbinMiddleware.RequireGroupPermission("chapter", "update"),
		cr.handler.updateChapter,
	)
}
