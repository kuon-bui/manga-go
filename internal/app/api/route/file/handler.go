package fileroute

import (
	"manga-go/internal/pkg/config"
	"manga-go/internal/pkg/logger"
	fileservice "manga-go/internal/pkg/services/file"

	"github.com/hibiken/asynq"
	"go.uber.org/fx"
)

type FileHandler struct {
	fileService *fileservice.FileService
	asynqClient *asynq.Client
	config      *config.Config
	logger      *logger.Logger
}

type FileHandlerParams struct {
	fx.In
	FileService *fileservice.FileService
	AsynqClient *asynq.Client
	Config      *config.Config
	Logger      *logger.Logger
}

func NewFileHandler(params FileHandlerParams) *FileHandler {
	return &FileHandler{
		fileService: params.FileService,
		asynqClient: params.AsynqClient,
		config:      params.Config,
		logger:      params.Logger,
	}
}
