package ratingservice

import (
	"context"
	"manga-go/internal/pkg/common"
	"manga-go/internal/pkg/logger"
	"manga-go/internal/pkg/model"
	comicrepo "manga-go/internal/pkg/repo/comic"
	ratingrepo "manga-go/internal/pkg/repo/rating"

	"github.com/google/uuid"
	"go.uber.org/fx"
)

// RatingRepository defines the data access interface for Rating.
type RatingRepository interface {
	Create(ctx context.Context, rating *model.Rating) error
	FindOne(ctx context.Context, conditions []any, moreKeys map[string]common.MoreKeyOption) (*model.Rating, error)
	Update(ctx context.Context, conditions []any, data map[string]any) error
	DeleteSoft(ctx context.Context, conditions []any) error
	FindPaginated(ctx context.Context, conditions []any, paging *common.Paging, moreKeys map[string]common.MoreKeyOption) ([]*model.Rating, int64, error)
	FindByUserAndComic(ctx context.Context, userID uuid.UUID, comicID uuid.UUID) (*model.Rating, error)
	GetAverageRatingByComicID(ctx context.Context, comicID uuid.UUID) (float64, int64, error)
}

// ComicRepository defines the subset of the ComicRepo interface used by RatingService.
type ComicRepository interface {
	FindOne(ctx context.Context, conditions []any, moreKeys map[string]common.MoreKeyOption) (*model.Comic, error)
}

type RatingService struct {
	logger     *logger.Logger
	ratingRepo RatingRepository
	comicRepo  ComicRepository
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

// NewRatingServiceWithRepo creates a RatingService with explicit repository interfaces,
// useful for unit testing.
func NewRatingServiceWithRepo(l *logger.Logger, ratingRepo RatingRepository) *RatingService {
	return &RatingService{
		logger:     l,
		ratingRepo: ratingRepo,
	}
}
