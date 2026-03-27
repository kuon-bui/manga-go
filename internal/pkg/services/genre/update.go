package genreservice

import (
	"context"
	"errors"
	"manga-go/internal/app/api/common/response"
	genrerequest "manga-go/internal/pkg/request/genre"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (s *GenreService) UpdateGenre(ctx context.Context, slug string, req *genrerequest.UpdateGenreRequest) response.Result {
	genre, err := s.genreRepo.FindOne(ctx, []any{
		clause.Eq{Column: "slug", Value: slug},
	}, nil)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.ResultNotFound("Genre")
		}
		s.logger.Error("Failed to find genre", "error", err)
		return response.ResultErrDb(err)
	}

	if err := s.genreRepo.Update(ctx, []any{
		clause.Eq{Column: "id", Value: genre.ID},
	}, map[string]any{
		"name":        req.Name,
		"slug":        req.Slug,
		"description": req.Description,
		"thumbnail":   req.Thumbnail,
	}); err != nil {
		s.logger.Error("Failed to update genre", "error", err)
		return response.ResultErrDb(err)
	}

	genre.Name = req.Name
	genre.Slug = req.Slug
	genre.Description = req.Description
	genre.Thumbnail = req.Thumbnail
	return response.ResultSuccess("Genre updated successfully", genre)
}
