package chapterserivce

import (
	"manga-go/internal/pkg/logger"
	"manga-go/internal/pkg/redis"
	chapterrepo "manga-go/internal/pkg/repo/chapter"
	comicrepo "manga-go/internal/pkg/repo/comic"
	readingprogressrepo "manga-go/internal/pkg/repo/reading_progress"
	usercomicreadrepo "manga-go/internal/pkg/repo/user_comic_read"

	"github.com/hibiken/asynq"
	"go.uber.org/fx"
)

type ChapterService struct {
	logger              *logger.Logger
	chapterRepo         *chapterrepo.ChapterRepo
	comicRepo           *comicrepo.ComicRepo
	readingProgressRepo *readingprogressrepo.ReadingProgressRepo
	userComicReadRepo   *usercomicreadrepo.UserComicReadRepo
	rds                 *redis.Redis
	asynqClient         *asynq.Client
}

type ChapterServiceParams struct {
	fx.In
	Logger              *logger.Logger
	ChapterRepo         *chapterrepo.ChapterRepo
	ComicRepo           *comicrepo.ComicRepo
	ReadingProgressRepo *readingprogressrepo.ReadingProgressRepo
	UserComicReadRepo   *usercomicreadrepo.UserComicReadRepo
	Redis               *redis.Redis
	AsynqClient         *asynq.Client
}

func NewChapterService(params ChapterServiceParams) *ChapterService {
	return &ChapterService{
		logger:              params.Logger,
		chapterRepo:         params.ChapterRepo,
		comicRepo:           params.ComicRepo,
		readingProgressRepo: params.ReadingProgressRepo,
		rds:                 params.Redis,
		userComicReadRepo:   params.UserComicReadRepo,
		asynqClient:         params.AsynqClient,
	}
}
