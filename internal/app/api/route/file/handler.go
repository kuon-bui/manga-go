package fileroute

import (
	fileservice "manga-go/internal/pkg/services/file"

	"go.uber.org/fx"
)

type FileHandler struct {
	fileService *fileservice.FileService
}

type FileHandlerParams struct {
	fx.In
	FileService *fileservice.FileService
}

func NewFileHandler(params FileHandlerParams) *FileHandler {
	return &FileHandler{
		fileService: params.FileService,
	}
}
