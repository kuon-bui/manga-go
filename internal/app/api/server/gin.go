package server

import (
	_ "manga-go/docs"
	"manga-go/internal/pkg/config"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

func NewGinEngine(config *config.Config) *gin.Engine {
	g := gin.Default()

	g.Use(gzip.Gzip(gzip.DefaultCompression))
	g.Use(otelgin.Middleware(config.Service.Name))
	g.Use(gin.Recovery())

	ginConfig := cors.DefaultConfig()
	ginConfig.AllowOrigins = strings.Split(config.Service.AllowOrigins, ",")
	ginConfig.AllowCredentials = true
	ginConfig.AllowHeaders = []string{
		"Access-Control-Allow-Origin",
		"Origin",
		"Content-Type",
		"Content-Length",
		"Accept-Encoding",
		"X-CSRF-Token",
		"Authorization",
		"Referer",
		"X-Size",
		"Credentials",
	}
	ginConfig.ExposeHeaders = []string{"Content-Disposition"}
	g.Use(cors.New(ginConfig))

	// Swagger API documentation route
	g.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return g
}
