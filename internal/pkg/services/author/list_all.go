package authorservice

import (
	"context"
	"manga-go/internal/app/api/common/response"
)

func (s *AuthorService) ListAllAuthors(ctx context.Context) response.Result {
	authors, err := s.authorRepo.FindAll(ctx, nil, nil)
	if err != nil {
		s.logger.Error("Failed to list all authors", "error", err)
		return response.ResultErrDb(err)
	}

	return response.ResultSuccess("Authors retrieved successfully", authors)
}
