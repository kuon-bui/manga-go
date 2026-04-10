package genreservice

import (
	"context"
	"manga-go/internal/pkg/common"
	"manga-go/internal/pkg/logger"
	"manga-go/internal/pkg/model"
	"manga-go/internal/pkg/redis"
	genrerepo "manga-go/internal/pkg/repo/genre"

	"go.uber.org/fx"
)

// GenreRepository defines the data access interface for Genre.
type GenreRepository interface {
	Create(ctx context.Context, genre *model.Genre) error
	FindOne(ctx context.Context, conditions []any, moreKeys map[string]common.MoreKeyOption) (*model.Genre, error)
	Update(ctx context.Context, conditions []any, data map[string]any) error
	DeleteSoft(ctx context.Context, conditions []any) error
	FindPaginated(ctx context.Context, conditions []any, paging *common.Paging, moreKeys map[string]common.MoreKeyOption) ([]*model.Genre, int64, error)
	FindAll(ctx context.Context, conditions []any, moreKeys map[string]common.MoreKeyOption) ([]*model.Genre, error)
}

type GenreService struct {
	logger    *logger.Logger
	genreRepo GenreRepository
	rds       *redis.Redis
}

type GenreServiceParams struct {
	fx.In
	Logger    *logger.Logger
	GenreRepo *genrerepo.GenreRepo
	Redis     *redis.Redis
}

func NewGenreService(params GenreServiceParams) *GenreService {
	return &GenreService{
		logger:    params.Logger,
		genreRepo: params.GenreRepo,
		rds:       params.Redis,
	}
}

// NewGenreServiceWithRepo creates a GenreService with an explicit repository,
// useful for unit testing.
func NewGenreServiceWithRepo(l *logger.Logger, repo GenreRepository) *GenreService {
	return &GenreService{
		logger:    l,
		genreRepo: repo,
	}
}
