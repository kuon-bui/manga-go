package authorservice

import (
	"context"
	"manga-go/internal/app/api/common/response"
	"manga-go/internal/pkg/model"
	authorrequest "manga-go/internal/pkg/request/author"
)

func (s *AuthorService) CreateAuthor(ctx context.Context, req *authorrequest.CreateAuthorRequest) response.Result {
	author := model.Author{
		Name: req.Name,
	}

	if err := s.authorRepo.Create(ctx, &author); err != nil {
		s.logger.Error("Failed to create author", "error", err)
		return response.ResultErrDb(err)
	}

	return response.ResultSuccess("Author created successfully", author)
}
