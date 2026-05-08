package ratingservice

import (
	"manga-go/internal/pkg/logger"
	comicrepo "manga-go/internal/pkg/repo/comic"
	ratingrepo "manga-go/internal/pkg/repo/rating"

	"github.com/hibiken/asynq"
	"go.uber.org/fx"
)

type RatingService struct {
	logger      *logger.Logger
	ratingRepo  *ratingrepo.RatingRepo
	comicRepo   *comicrepo.ComicRepo
	asynqClient *asynq.Client
}

type RatingServiceParams struct {
	fx.In
	Logger      *logger.Logger
	RatingRepo  *ratingrepo.RatingRepo
	ComicRepo   *comicrepo.ComicRepo
	AsynqClient *asynq.Client
}

func NewRatingService(p RatingServiceParams) *RatingService {
	return &RatingService{
		logger:      p.Logger,
		ratingRepo:  p.RatingRepo,
		comicRepo:   p.ComicRepo,
		asynqClient: p.AsynqClient,
	}
}
