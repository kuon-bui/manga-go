package genreservice

import (
	"context"
	"manga-go/internal/app/api/common/response"
	"manga-go/internal/pkg/model"
	genrerequest "manga-go/internal/pkg/request/genre"
)

func (s *GenreService) CreateGenre(ctx context.Context, req *genrerequest.CreateGenreRequest) response.Result {
	genre := model.Genre{
		Name:        req.Name,
		Slug:        req.Slug,
		Description: req.Description,
		Thumbnail:   req.Thumbnail,
	}

	if err := s.genreRepo.Create(ctx, &genre); err != nil {
		s.logger.Error("Failed to create genre", "error", err)
		return response.ResultErrDb(err)
	}

	return response.ResultSuccess("Genre created successfully", genre)
}
