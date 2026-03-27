package chapterhandler

import (
	"manga-go/internal/pkg/logger"
	chapterserivce "manga-go/internal/pkg/services/chapter"

	"go.uber.org/fx"
)

type ChapterHandler struct {
	logger         *logger.Logger
	chapterService *chapterserivce.ChapterService
}

type ChapterHandlerParams struct {
	fx.In
	Logger         *logger.Logger
	ChapterService *chapterserivce.ChapterService
}

func NewChapterHandler(params ChapterHandlerParams) *ChapterHandler {
	return &ChapterHandler{
		logger:         params.Logger,
		chapterService: params.ChapterService,
	}
}
