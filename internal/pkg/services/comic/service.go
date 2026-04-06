package comicservice

import (
	"manga-go/internal/pkg/logger"
	"manga-go/internal/pkg/redis"
	comicrepo "manga-go/internal/pkg/repo/comic"
	genrerepo "manga-go/internal/pkg/repo/genre"
	tagrepo "manga-go/internal/pkg/repo/tag"
	usercomicreadrepo "manga-go/internal/pkg/repo/user_comic_read"

	"go.uber.org/fx"
)

type ComicService struct {
	logger            *logger.Logger
	comicRepo         *comicrepo.ComicRepo
	genreRepo         *genrerepo.GenreRepo
	tagRepo           *tagrepo.TagRepo
	userComicReadRepo *usercomicreadrepo.UserComicReadRepo
	rds               *redis.Redis
}

type ComicServiceParams struct {
	fx.In
	Logger            *logger.Logger
	ComicRepo         *comicrepo.ComicRepo
	GenreRepo         *genrerepo.GenreRepo
	TagRepo           *tagrepo.TagRepo
	UserComicReadRepo *usercomicreadrepo.UserComicReadRepo
	Redis             *redis.Redis
}

func NewComicService(params ComicServiceParams) *ComicService {
	return &ComicService{
		logger:            params.Logger,
		comicRepo:         params.ComicRepo,
		genreRepo:         params.GenreRepo,
		tagRepo:           params.TagRepo,
		userComicReadRepo: params.UserComicReadRepo,
		rds:               params.Redis,
	}
}
