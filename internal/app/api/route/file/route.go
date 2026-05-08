package fileroute

import (
	authmiddleware "manga-go/internal/app/middleware/auth"
	authzmiddleware "manga-go/internal/app/middleware/authz"
	"manga-go/internal/pkg/authorization"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

type FileRoute struct {
	*gin.Engine
	fileHandler     *FileHandler
	authmiddleware  *authmiddleware.AuthMiddleware
	authzMiddleware *authzmiddleware.AuthzMiddleware
}

type FileRouteParams struct {
	fx.In
	R               *gin.Engine
	FileHandler     *FileHandler
	AuthMiddleware  *authmiddleware.AuthMiddleware
	AuthzMiddleware *authzmiddleware.AuthzMiddleware
}

func NewFileRoute(params FileRouteParams) *FileRoute {
	return &FileRoute{
		Engine:          params.R,
		fileHandler:     params.FileHandler,
		authmiddleware:  params.AuthMiddleware,
		authzMiddleware: params.AuthzMiddleware,
	}
}

func (fr *FileRoute) Setup() {
	rg := fr.Group("/files")
	privateRg := rg.Group("", fr.authmiddleware.RequireJwt)
	requireFileCreate := authzmiddleware.Require(fr.authzMiddleware, authorization.ActionCreate, authorization.ObjectFile)

	privateRg.POST("/upload", requireFileCreate, fr.fileHandler.uploadImage)
	privateRg.GET("/presign/*filename", requireFileCreate, fr.fileHandler.getPresignURL)
	rg.GET("/content/*filename", fr.fileHandler.getFileContent)
}
