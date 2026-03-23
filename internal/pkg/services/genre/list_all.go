package genreservice

import (
	"context"
	"manga-go/internal/app/api/common/response"
)

func (s *GenreService) ListAllGenres(ctx context.Context) response.Result {
	genres, err := s.genreRepo.FindAll(ctx, nil, nil)
	if err != nil {
		s.logger.Error("Failed to list all genres", "error", err)
		return response.ResultErrDb(err)
	}

	return response.ResultSuccess("Genres retrieved successfully", genres)
}
