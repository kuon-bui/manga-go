package genreroute

import (
	authmiddleware "manga-go/internal/app/middleware/auth"
	"manga-go/internal/pkg/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

type GenreRoute struct {
	logger         *logger.Logger
	r              *gin.Engine
	authMiddleware *authmiddleware.AuthMiddleware
	genreHandler   *GenreHandler
}

type GenreRouteParams struct {
	fx.In

	R              *gin.Engine
	Logger         *logger.Logger
	GenreHandler   *GenreHandler
	AuthMiddleware *authmiddleware.AuthMiddleware
}

func NewGenreRoute(params GenreRouteParams) *GenreRoute {
	return &GenreRoute{
		logger:         params.Logger,
		r:              params.R,
		authMiddleware: params.AuthMiddleware,
		genreHandler:   params.GenreHandler,
	}
}

func (gr *GenreRoute) Setup() {
	rg := gr.r.Group("/genres", gr.authMiddleware.RequireJwt)

	rg.GET("/", gr.genreHandler.getGenres)
	rg.GET("/all", gr.genreHandler.getAllGenres)
	rg.GET("/:id", gr.genreHandler.getGenre)
	rg.POST("/", gr.genreHandler.createGenre)
	rg.PUT("/:id", gr.genreHandler.updateGenre)
	rg.DELETE("/:id", gr.genreHandler.deleteGenre)
}
