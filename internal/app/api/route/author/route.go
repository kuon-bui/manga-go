package authorroute

import (
	authmiddleware "manga-go/internal/app/middleware/auth"
	authzmiddleware "manga-go/internal/app/middleware/authz"
	"manga-go/internal/pkg/authorization"
	"manga-go/internal/pkg/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

type AuthorRoute struct {
	*gin.Engine
	logger          *logger.Logger
	authMiddleware  *authmiddleware.AuthMiddleware
	authzMiddleware *authzmiddleware.AuthzMiddleware
	authorHandler   *AuthorHandler
}

type AuthorRouteParams struct {
	fx.In

	R               *gin.Engine
	Logger          *logger.Logger
	AuthorHandler   *AuthorHandler
	AuthMiddleware  *authmiddleware.AuthMiddleware
	AuthzMiddleware *authzmiddleware.AuthzMiddleware
}

func NewAuthorRoute(params AuthorRouteParams) *AuthorRoute {
	return &AuthorRoute{
		Engine:          params.R,
		logger:          params.Logger,
		authMiddleware:  params.AuthMiddleware,
		authzMiddleware: params.AuthzMiddleware,
		authorHandler:   params.AuthorHandler,
	}
}

func (ar *AuthorRoute) Setup() {
	rg := ar.Group("/authors", ar.authMiddleware.RequireJwt)
	requireAuthorWrite := authzmiddleware.Require(ar.authzMiddleware, authorization.ActionCreate, authorization.ObjectAuthor)
	requireAuthorUpdate := authzmiddleware.Require(ar.authzMiddleware, authorization.ActionUpdate, authorization.ObjectAuthor)
	requireAuthorDelete := authzmiddleware.Require(ar.authzMiddleware, authorization.ActionDelete, authorization.ObjectAuthor)

	rg.GET("", ar.authorHandler.getAuthors)
	rg.GET("/all", ar.authorHandler.getAllAuthors)
	rg.GET("/:id", ar.authorHandler.getAuthor)
	rg.POST("", requireAuthorWrite, ar.authorHandler.createAuthor)
	rg.PUT("/:id", requireAuthorUpdate, ar.authorHandler.updateAuthor)
	rg.DELETE("/:id", requireAuthorDelete, ar.authorHandler.deleteAuthor)
}
