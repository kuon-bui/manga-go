package server

import (
	"manga-go/internal/pkg/config"
	validatorpkg "manga-go/internal/pkg/validator"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

func NewGinEngine(config *config.Config) *gin.Engine {
	gin.ForceConsoleColor()
	g := gin.Default()
	g.Use(gzip.Gzip(gzip.DefaultCompression))
	g.Use(otelgin.Middleware(config.Service.Name))
	g.Use(gin.Recovery())

	ginConfig := cors.DefaultConfig()

	if config.Service.AllowOrigins != "" && config.Service.AllowOrigins != "*" {
		ginConfig.AllowOrigins = strings.Split(config.Service.AllowOrigins, ",")
	}

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

	validatorFuncs := map[string]validator.Func{
		"age_rating": validatorpkg.ValidateAgeRating,
		"comic_type": validatorpkg.ValidateComicType,
	}

	for tag, fn := range validatorFuncs {
		if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
			v.RegisterValidation(tag, fn)
		}
	}

	ginConfig.ExposeHeaders = []string{"Content-Disposition"}
	g.Use(cors.New(ginConfig))

	return g
}
