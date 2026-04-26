package image

import (
	"manga-go/internal/pkg/logger"
	fileservice "manga-go/internal/pkg/services/file"
	queueconstant "manga-go/internal/queue/queue_constant"

	"github.com/hibiken/asynq"
	"go.uber.org/fx"
)

type ImageProcessWorkerParams struct {
	fx.In
	Mux         *asynq.ServeMux
	Logger      *logger.Logger
	FileService *fileservice.FileService
}

type ImageProcessWorker struct {
	mux         *asynq.ServeMux
	logger      *logger.Logger
	fileService *fileservice.FileService
}

func NewImageProcessWorker(p ImageProcessWorkerParams) *ImageProcessWorker {
	return &ImageProcessWorker{
		mux:         p.Mux,
		logger:      p.Logger,
		fileService: p.FileService,
	}
}

func (w *ImageProcessWorker) RegisterWorkers() {
	w.mux.HandleFunc(queueconstant.IMAGE_PROCESS_TASK, w.imageProcessHandler)
	w.mux.HandleFunc(queueconstant.IMAGE_PROCESS_CLEANUP_TASK, w.imageProcessCleanupHandler)
}
