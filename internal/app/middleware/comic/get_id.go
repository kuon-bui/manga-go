package comicmiddleware

import (
	"manga-go/internal/pkg/common"
	"manga-go/internal/pkg/logger"
	comicservice "manga-go/internal/pkg/services/comic"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

type ComicMiddleware struct {
	logger       *logger.Logger
	comicService *comicservice.ComicService
}

type ComicMiddlewareParams struct {
	fx.In
	Logger       *logger.Logger
	ComicService *comicservice.ComicService
}

func NewComicMiddleware(params ComicMiddlewareParams) *ComicMiddleware {
	return &ComicMiddleware{
		logger:       params.Logger,
		comicService: params.ComicService,
	}
}

func (m *ComicMiddleware) ResolveComicID(g *gin.Context) {
	comicSlug := g.Param("comicSlug")
	if comicSlug == "" {
		g.AbortWithStatusJSON(400, gin.H{"error": "comicSlug parameter is required"})
		return
	}

	id, err := m.comicService.GetComicIDBySlug(g.Request.Context(), comicSlug)
	if err != nil {
		m.logger.Errorf("Failed to get comic ID for slug %s: %v", comicSlug, err)
		g.AbortWithStatusJSON(404, gin.H{"error": "Comic not found"})
		return
	}

	ctx := common.SetComicIdToContext(g.Request.Context(), id)
	g.Request = g.Request.WithContext(ctx)
}
