package authorroute

import (
	authmiddleware "manga-go/internal/app/middleware/auth"
	"manga-go/internal/pkg/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

type AuthorRoute struct {
	*gin.Engine
	logger         *logger.Logger
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
		Engine:         params.R,
		logger:         params.Logger,
		authMiddleware: params.AuthMiddleware,
		authorHandler:  params.AuthorHandler,
	}
}

func (ar *AuthorRoute) Setup() {
	rg := ar.Group("/authors", ar.authMiddleware.RequireJwt)

	rg.GET("", ar.authorHandler.getAuthors)
	rg.GET("/all", ar.authorHandler.getAllAuthors)
	rg.GET("/:id", ar.authorHandler.getAuthor)
	rg.POST("", ar.authorHandler.createAuthor)
	rg.PUT("/:id", ar.authorHandler.updateAuthor)
	rg.DELETE("/:id", ar.authorHandler.deleteAuthor)
}
