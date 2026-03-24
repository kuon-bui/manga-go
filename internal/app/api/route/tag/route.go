package tagroute

import (
	authmiddleware "manga-go/internal/app/middleware/auth"
	"manga-go/internal/pkg/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

type TagRoute struct {
	logger         *logger.Logger
	r              *gin.Engine
	authMiddleware *authmiddleware.AuthMiddleware
	tagHandler     *TagHandler
}

type TagRouteParams struct {
	fx.In

	R              *gin.Engine
	Logger         *logger.Logger
	TagHandler     *TagHandler
	AuthMiddleware *authmiddleware.AuthMiddleware
}

func NewTagRoute(params TagRouteParams) *TagRoute {
	return &TagRoute{
		logger:         params.Logger,
		r:              params.R,
		authMiddleware: params.AuthMiddleware,
		tagHandler:     params.TagHandler,
	}
}

func (tr *TagRoute) Setup() {
	rg := tr.r.Group("/tags", tr.authMiddleware.RequireJwt)

	rg.GET("", tr.tagHandler.getTags)
	rg.GET("/all", tr.tagHandler.getAllTags)
	rg.POST("", tr.tagHandler.createTag)
	rg.DELETE("/:slug", tr.tagHandler.deleteTag)
}
