package comicservice

import (
	"manga-go/internal/pkg/logger"
	"manga-go/internal/pkg/redis"
	authorrepo "manga-go/internal/pkg/repo/author"
	comicrepo "manga-go/internal/pkg/repo/comic"
	comicfollowrepo "manga-go/internal/pkg/repo/comic_follow"
	genrerepo "manga-go/internal/pkg/repo/genre"
	tagrepo "manga-go/internal/pkg/repo/tag"
	translationgrouprepo "manga-go/internal/pkg/repo/translation_group"
	usercomicreadrepo "manga-go/internal/pkg/repo/user_comic_read"

	"go.uber.org/fx"
	"gorm.io/gorm"
)

type ComicService struct {
	logger               *logger.Logger
	authorRepo           *authorrepo.AuthorRepo
	comicFollowRepo      *comicfollowrepo.ComicFollowRepo
	comicRepo            *comicrepo.ComicRepo
	genreRepo            *genrerepo.GenreRepo
	tagRepo              *tagrepo.TagRepo
	translationGroupRepo *translationgrouprepo.TranslationGroupRepo
	userComicReadRepo    *usercomicreadrepo.UserComicReadRepo
	rds                  *redis.Redis
	gormDb               *gorm.DB
}

type ComicServiceParams struct {
	fx.In
	Logger               *logger.Logger
	AuthorRepo           *authorrepo.AuthorRepo
	ComicFollowRepo      *comicfollowrepo.ComicFollowRepo
	ComicRepo            *comicrepo.ComicRepo
	GenreRepo            *genrerepo.GenreRepo
	TagRepo              *tagrepo.TagRepo
	TranslationGroupRepo *translationgrouprepo.TranslationGroupRepo
	UserComicReadRepo    *usercomicreadrepo.UserComicReadRepo
	Redis                *redis.Redis
	GormDb               *gorm.DB
}

func NewComicService(params ComicServiceParams) *ComicService {
	return &ComicService{
		logger:               params.Logger,
		authorRepo:           params.AuthorRepo,
		comicFollowRepo:      params.ComicFollowRepo,
		comicRepo:            params.ComicRepo,
		genreRepo:            params.GenreRepo,
		tagRepo:              params.TagRepo,
		translationGroupRepo: params.TranslationGroupRepo,
		userComicReadRepo:    params.UserComicReadRepo,
		rds:                  params.Redis,
		gormDb:               params.GormDb,
	}
}
