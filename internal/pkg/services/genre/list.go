package genreservice

import (
	"context"
	"manga-go/internal/app/api/common/response"
	"manga-go/internal/pkg/common"
)

func (s *GenreService) ListGenres(ctx context.Context, paging *common.Paging) response.Result {
	genres, total, err := s.genreRepo.FindPaginated(ctx, nil, paging, nil)
	if err != nil {
		s.logger.Error("Failed to list genres", "error", err)
		return response.ResultErrDb(err)
	}

	return response.ResultSuccess("Genres retrieved successfully", response.ResponsePaginationData(genres, total))
}
