package fileroute

import (
	authmiddleware "manga-go/internal/app/middleware/auth"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

type FileRoute struct {
	*gin.Engine
	fileHandler    *FileHandler
	authmiddleware *authmiddleware.AuthMiddleware
}

type FileRouteParams struct {
	fx.In
	R              *gin.Engine
	FileHandler    *FileHandler
	AuthMiddleware *authmiddleware.AuthMiddleware
}

func NewFileRoute(params FileRouteParams) *FileRoute {
	return &FileRoute{
		Engine:         params.R,
		fileHandler:    params.FileHandler,
		authmiddleware: params.AuthMiddleware,
	}
}

func (fr *FileRoute) Setup() {
	rg := fr.Group("/files")
	privateRg := rg.Group("", fr.authmiddleware.RequireJwt)

	privateRg.POST("/upload", fr.fileHandler.uploadImage)
	privateRg.GET("/presign/*filename", fr.fileHandler.getPresignURL)
	rg.GET("/content/*filename", fr.fileHandler.getFileContent)
}
