package tagroute

import (
	authmiddleware "manga-go/internal/app/middleware/auth"
	authzmiddleware "manga-go/internal/app/middleware/authz"
	"manga-go/internal/pkg/authorization"
	"manga-go/internal/pkg/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

type TagRoute struct {
	*gin.Engine
	logger          *logger.Logger
	authMiddleware  *authmiddleware.AuthMiddleware
	authzMiddleware *authzmiddleware.AuthzMiddleware
	tagHandler      *TagHandler
}

type TagRouteParams struct {
	fx.In

	R               *gin.Engine
	Logger          *logger.Logger
	TagHandler      *TagHandler
	AuthMiddleware  *authmiddleware.AuthMiddleware
	AuthzMiddleware *authzmiddleware.AuthzMiddleware
}

func NewTagRoute(params TagRouteParams) *TagRoute {
	return &TagRoute{
		Engine:          params.R,
		logger:          params.Logger,
		authMiddleware:  params.AuthMiddleware,
		authzMiddleware: params.AuthzMiddleware,
		tagHandler:      params.TagHandler,
	}
}

func (tr *TagRoute) Setup() {
	rg := tr.Group("/tags", tr.authMiddleware.RequireJwt)
	requireTagWrite := authzmiddleware.Require(tr.authzMiddleware, authorization.ActionCreate, authorization.ObjectTag)
	requireTagDelete := authzmiddleware.Require(tr.authzMiddleware, authorization.ActionDelete, authorization.ObjectTag)

	rg.GET("", tr.tagHandler.getTags)
	rg.GET("/all", tr.tagHandler.getAllTags)
	rg.POST("", requireTagWrite, tr.tagHandler.createTag)
	rg.DELETE("/:slug", requireTagDelete, tr.tagHandler.deleteTag)
}
