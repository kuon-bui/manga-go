package authorservice

import (
	"context"
	"manga-go/internal/pkg/common"
	"manga-go/internal/pkg/logger"
	"manga-go/internal/pkg/model"
	authorrepo "manga-go/internal/pkg/repo/author"

	"go.uber.org/fx"
)

// AuthorRepository defines the data access interface for Author.
type AuthorRepository interface {
	Create(ctx context.Context, author *model.Author) error
	FindOne(ctx context.Context, conditions []any, moreKeys map[string]common.MoreKeyOption) (*model.Author, error)
	Update(ctx context.Context, conditions []any, data map[string]any) error
	DeleteSoft(ctx context.Context, conditions []any) error
	FindPaginated(ctx context.Context, conditions []any, paging *common.Paging, moreKeys map[string]common.MoreKeyOption) ([]*model.Author, int64, error)
	FindAll(ctx context.Context, conditions []any, moreKeys map[string]common.MoreKeyOption) ([]*model.Author, error)
}

type AuthorService struct {
	logger     *logger.Logger
	authorRepo AuthorRepository
}

type AuthorServiceParams struct {
	fx.In
	Logger     *logger.Logger
	AuthorRepo *authorrepo.AuthorRepo
}

func NewAuthorService(params AuthorServiceParams) *AuthorService {
	return &AuthorService{
		logger:     params.Logger,
		authorRepo: params.AuthorRepo,
	}
}

// NewAuthorServiceWithRepo creates an AuthorService with an explicit repository,
// useful for unit testing.
func NewAuthorServiceWithRepo(l *logger.Logger, repo AuthorRepository) *AuthorService {
	return &AuthorService{
		logger:     l,
		authorRepo: repo,
	}
}
