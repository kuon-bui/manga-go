package chapterserivce

import (
	"manga-go/internal/pkg/logger"
	chapterrepo "manga-go/internal/pkg/repo/chapter"
	comicrepo "manga-go/internal/pkg/repo/comic"

	"go.uber.org/fx"
)

type ChapterService struct {
	logger      *logger.Logger
	chapterRepo *chapterrepo.ChapterRepo
	comicRepo   *comicrepo.ComicRepo
}

type ChapterServiceParams struct {
	fx.In
	Logger      *logger.Logger
	ChapterRepo *chapterrepo.ChapterRepo
	ComicRepo   *comicrepo.ComicRepo
}

func NewChapterService(params ChapterServiceParams) *ChapterService {
	return &ChapterService{
		logger:      params.Logger,
		chapterRepo: params.ChapterRepo,
		comicRepo:   params.ComicRepo,
	}
}
