package ratingservice

import (
	"manga-go/internal/pkg/logger"
	comicrepo "manga-go/internal/pkg/repo/comic"
	ratingrepo "manga-go/internal/pkg/repo/rating"

	"go.uber.org/fx"
)

type RatingService struct {
	logger     *logger.Logger
	ratingRepo *ratingrepo.RatingRepo
	comicRepo  *comicrepo.ComicRepo
}

type RatingServiceParams struct {
	fx.In
	Logger     *logger.Logger
	RatingRepo *ratingrepo.RatingRepo
	ComicRepo  *comicrepo.ComicRepo
}

func NewRatingService(p RatingServiceParams) *RatingService {
	return &RatingService{
		logger:     p.Logger,
		ratingRepo: p.RatingRepo,
		comicRepo:  p.ComicRepo,
	}
}
