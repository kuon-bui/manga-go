package authorservice

import (
	"manga-go/internal/pkg/logger"
	authorrepo "manga-go/internal/pkg/repo/author"

	"go.uber.org/fx"
)

type AuthorService struct {
	logger     *logger.Logger
	authorRepo *authorrepo.AuthorRepo
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
