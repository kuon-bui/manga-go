package authorroute

import (
	authorservice "manga-go/internal/pkg/services/author"

	"go.uber.org/fx"
)

type AuthorHandler struct {
	authorService *authorservice.AuthorService
}

type AuthorHandlerParams struct {
	fx.In

	AuthorService *authorservice.AuthorService
}

func NewAuthorHandler(p AuthorHandlerParams) *AuthorHandler {
	return &AuthorHandler{
		authorService: p.AuthorService,
	}
}
