package slugmiddleware

import (
	"manga-go/internal/pkg/common"
	"manga-go/internal/pkg/logger"
	chapterservice "manga-go/internal/pkg/services/chapter"
	comicservice "manga-go/internal/pkg/services/comic"
	genreservice "manga-go/internal/pkg/services/genre"
	tagservice "manga-go/internal/pkg/services/tag"
	translationgroupservice "manga-go/internal/pkg/services/translation_group"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

type SlugMiddleware struct {
	logger                  *logger.Logger
	comicService            *comicservice.ComicService
	chapterService          *chapterservice.ChapterService
	translationGroupService *translationgroupservice.TranslationGroupService
	genreService            *genreservice.GenreService
	tagService              *tagservice.TagService
}

type SlugMiddlewareParams struct {
	fx.In
	Logger                  *logger.Logger
	ComicService            *comicservice.ComicService
	ChapterService          *chapterservice.ChapterService
	TranslationGroupService *translationgroupservice.TranslationGroupService
	GenreService            *genreservice.GenreService
	TagService              *tagservice.TagService
}

func NewSlugMiddleware(params SlugMiddlewareParams) *SlugMiddleware {
	return &SlugMiddleware{
		logger:                  params.Logger,
		comicService:            params.ComicService,
		chapterService:          params.ChapterService,
		translationGroupService: params.TranslationGroupService,
		genreService:            params.GenreService,
		tagService:              params.TagService,
	}
}

func (m *SlugMiddleware) ResolveComicID(g *gin.Context) {
	comicSlug := g.Param("comicSlug")
	if comicSlug == "" {
		g.AbortWithStatusJSON(400, gin.H{"error": "comicSlug parameter is required"})
		return
	}

	id, groupID, err := m.comicService.GetComicIDAndGroupIDBySlug(g.Request.Context(), comicSlug)
	if err != nil {
		m.logger.Errorf("Failed to get comic ID for slug %s: %v", comicSlug, err)
		g.AbortWithStatusJSON(404, gin.H{"error": "Comic not found"})
		return
	}

	ctx := common.SetComicIdToContext(g.Request.Context(), id)
	if groupID != nil {
		ctx = common.SetTranslationGroupIdToContext(ctx, *groupID)
	}
	g.Request = g.Request.WithContext(ctx)
}

func (m *SlugMiddleware) ResolveChapterID(g *gin.Context) {
	chapterSlug := g.Param("chapterSlug")
	if chapterSlug == "" {
		g.AbortWithStatusJSON(400, gin.H{"error": "chapterSlug parameter is required"})
		return
	}

	id, err := m.chapterService.GetChapterIDBySlug(g.Request.Context(), chapterSlug)
	if err != nil {
		m.logger.Errorf("Failed to get chapter ID for slug %s: %v", chapterSlug, err)
		g.AbortWithStatusJSON(404, gin.H{"error": "Chapter not found"})
		return
	}

	ctx := common.SetChapterIdToContext(g.Request.Context(), id)
	g.Request = g.Request.WithContext(ctx)
}

func (m *SlugMiddleware) ResolveTranslationGroupID(g *gin.Context) {
	slug := g.Param("slug")
	if slug == "" {
		g.AbortWithStatusJSON(400, gin.H{"error": "slug parameter is required"})
		return
	}

	id, err := m.translationGroupService.GetTranslationGroupIDBySlug(g.Request.Context(), slug)
	if err != nil {
		m.logger.Errorf("Failed to get translation group ID for slug %s: %v", slug, err)
		g.AbortWithStatusJSON(404, gin.H{"error": "Translation group not found"})
		return
	}

	ctx := common.SetTranslationGroupIdToContext(g.Request.Context(), id)
	g.Request = g.Request.WithContext(ctx)
}

func (m *SlugMiddleware) ResolveGenreID(g *gin.Context) {
	genreSlug := g.Param("genreSlug")
	if genreSlug == "" {
		g.AbortWithStatusJSON(400, gin.H{"error": "genreSlug parameter is required"})
		return
	}

	id, err := m.genreService.GetGenreIDBySlug(g.Request.Context(), genreSlug)
	if err != nil {
		m.logger.Errorf("Failed to get genre ID for slug %s: %v", genreSlug, err)
		g.AbortWithStatusJSON(404, gin.H{"error": "Genre not found"})
		return
	}

	ctx := common.SetGenreIdToContext(g.Request.Context(), id)
	g.Request = g.Request.WithContext(ctx)
}

func (m *SlugMiddleware) ResolveTagID(g *gin.Context) {
	tagSlug := g.Param("tagSlug")
	if tagSlug == "" {
		g.AbortWithStatusJSON(400, gin.H{"error": "tagSlug parameter is required"})
		return
	}

	id, err := m.tagService.GetTagIDBySlug(g.Request.Context(), tagSlug)
	if err != nil {
		m.logger.Errorf("Failed to get tag ID for slug %s: %v", tagSlug, err)
		g.AbortWithStatusJSON(404, gin.H{"error": "Tag not found"})
		return
	}

	ctx := common.SetTagIdToContext(g.Request.Context(), id)
	g.Request = g.Request.WithContext(ctx)
}
