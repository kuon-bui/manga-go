package authorroute

import (
	authmiddleware "manga-go/internal/app/middleware/auth"
	"manga-go/internal/pkg/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

type AuthorRoute struct {
	logger         *logger.Logger
	r              *gin.Engine
	authMiddleware *authmiddleware.AuthMiddleware
	authorHandler  *AuthorHandler
}

type AuthorRouteParams struct {
	fx.In

	R              *gin.Engine
	Logger         *logger.Logger
	AuthorHandler  *AuthorHandler
	AuthMiddleware *authmiddleware.AuthMiddleware
}

func NewAuthorRoute(params AuthorRouteParams) *AuthorRoute {
	return &AuthorRoute{
		logger:         params.Logger,
		r:              params.R,
		authMiddleware: params.AuthMiddleware,
		authorHandler:  params.AuthorHandler,
	}
}

func (ar *AuthorRoute) Setup() {
	rg := ar.r.Group("/authors", ar.authMiddleware.RequireJwt)

	rg.GET("", ar.authorHandler.getAuthors)
	rg.GET("/all", ar.authorHandler.getAllAuthors)
	rg.GET("/:id", ar.authorHandler.getAuthor)
	rg.POST("", ar.authorHandler.createAuthor)
	rg.PUT("/:id", ar.authorHandler.updateAuthor)
	rg.DELETE("/:id", ar.authorHandler.deleteAuthor)
}
