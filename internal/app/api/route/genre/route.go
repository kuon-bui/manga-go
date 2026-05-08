package genreroute

import (
	authmiddleware "manga-go/internal/app/middleware/auth"
	authzmiddleware "manga-go/internal/app/middleware/authz"
	"manga-go/internal/pkg/authorization"
	"manga-go/internal/pkg/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

type GenreRoute struct {
	logger *logger.Logger
	*gin.Engine
	authMiddleware  *authmiddleware.AuthMiddleware
	authzMiddleware *authzmiddleware.AuthzMiddleware
	genreHandler    *GenreHandler
}

type GenreRouteParams struct {
	fx.In

	R               *gin.Engine
	Logger          *logger.Logger
	GenreHandler    *GenreHandler
	AuthMiddleware  *authmiddleware.AuthMiddleware
	AuthzMiddleware *authzmiddleware.AuthzMiddleware
}

func NewGenreRoute(params GenreRouteParams) *GenreRoute {
	return &GenreRoute{
		logger:          params.Logger,
		Engine:          params.R,
		authMiddleware:  params.AuthMiddleware,
		authzMiddleware: params.AuthzMiddleware,
		genreHandler:    params.GenreHandler,
	}
}

func (gr *GenreRoute) Setup() {
	rg := gr.Group("/genres", gr.authMiddleware.RequireJwt)
	requireGenreWrite := authzmiddleware.Require(gr.authzMiddleware, authorization.ActionCreate, authorization.ObjectGenre)
	requireGenreUpdate := authzmiddleware.Require(gr.authzMiddleware, authorization.ActionUpdate, authorization.ObjectGenre)
	requireGenreDelete := authzmiddleware.Require(gr.authzMiddleware, authorization.ActionDelete, authorization.ObjectGenre)

	rg.GET("", gr.genreHandler.getGenres)
	rg.GET("/all", gr.genreHandler.getAllGenres)
	rg.GET("/:slug", gr.genreHandler.getGenre)
	rg.POST("", requireGenreWrite, gr.genreHandler.createGenre)
	rg.PUT("/:slug", requireGenreUpdate, gr.genreHandler.updateGenre)
	rg.DELETE("/:slug", requireGenreDelete, gr.genreHandler.deleteGenre)
}
