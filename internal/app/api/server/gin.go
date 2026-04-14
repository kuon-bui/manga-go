package server

import (
	"manga-go/internal/pkg/config"
	validatorpkg "manga-go/internal/pkg/validator"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

func buildCorsConfig(cfg *config.Config) cors.Config {
	ginConfig := cors.Config{
		AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders: []string{
			"Origin",
			"Content-Type",
			"Content-Length",
			"Accept",
			"Accept-Encoding",
			"Authorization",
			"X-CSRF-Token",
		},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}

	rawOrigins := strings.TrimSpace(cfg.Service.AllowOrigins)
	if rawOrigins == "" {
		// In local development, allow dynamic localhost origins while keeping credentials enabled.
		ginConfig.AllowOriginFunc = func(origin string) bool {
			return strings.HasPrefix(origin, "http://localhost:") ||
				strings.HasPrefix(origin, "https://localhost:") ||
				strings.HasPrefix(origin, "http://127.0.0.1:") ||
				strings.HasPrefix(origin, "https://127.0.0.1:")
		}
		return ginConfig
	}

	if rawOrigins == "*" {
		// Keep wildcard behavior without using '*' response header, which is invalid with credentials.
		ginConfig.AllowOriginFunc = func(origin string) bool { return origin != "" }
		return ginConfig
	}

	origins := make([]string, 0)
	for origin := range strings.SplitSeq(rawOrigins, ",") {
		origin = strings.TrimSpace(origin)
		if origin == "" {
			continue
		}

		origins = append(origins, strings.TrimRight(origin, "/"))
	}

	if len(origins) == 0 {
		ginConfig.AllowOriginFunc = func(origin string) bool { return false }
		return ginConfig
	}

	ginConfig.AllowOrigins = origins
	return ginConfig
}

func NewGinEngine(config *config.Config) *gin.Engine {
	gin.ForceConsoleColor()
	g := gin.Default()

	validatorFuncs := map[string]validator.Func{
		"age_rating": validatorpkg.ValidateAgeRating,
		"comic_type": validatorpkg.ValidateComicType,
	}

	for tag, fn := range validatorFuncs {
		if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
			v.RegisterValidation(tag, fn)
		}
	}

	g.Use(cors.New(buildCorsConfig(config)))
	g.Use(gzip.Gzip(gzip.DefaultCompression))
	g.Use(otelgin.Middleware(config.Service.Name))
	g.Use(gin.Recovery())

	return g
}
