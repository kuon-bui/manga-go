package swaggerrouter

import (
	_ "manga-go/docs"
	authmiddleware "manga-go/internal/app/middleware/auth"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/fx"
)

type SwaggerRoute struct {
	r              *gin.Engine
	authMiddleware *authmiddleware.AuthMiddleware
}

type SwaggerRouteParams struct {
	fx.In

	R              *gin.Engine
	AuthMiddleware *authmiddleware.AuthMiddleware
}

func NewSwaggerRoute(params SwaggerRouteParams) *SwaggerRoute {
	return &SwaggerRoute{
		r:              params.R,
		authMiddleware: params.AuthMiddleware,
	}
}

func (sr *SwaggerRoute) Setup() {
	sr.r.GET(
		"/swagger/*any",
		ginSwagger.WrapHandler(
			swaggerFiles.Handler,
			func(c *ginSwagger.Config) {
				c.Title = "MangaGo API Documentation"
			},
		),
	)
}
