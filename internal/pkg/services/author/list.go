package authorservice

import (
	"context"
	"manga-go/internal/app/api/common/response"
	"manga-go/internal/pkg/common"
)

func (s *AuthorService) ListAuthors(ctx context.Context, paging *common.Paging) response.Result {
	authors, total, err := s.authorRepo.FindPaginated(ctx, nil, paging, nil)
	if err != nil {
		s.logger.Error("Failed to list authors", "error", err)
		return response.ResultErrDb(err)
	}

	return response.ResultSuccess("Authors retrieved successfully", response.ResponsePaginationData(authors, total))
}
