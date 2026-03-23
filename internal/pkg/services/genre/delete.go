package genreservice

import (
	"context"
	"errors"
	"manga-go/internal/app/api/common/response"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (s *GenreService) DeleteGenre(ctx context.Context, slug string) response.Result {
	genre, err := s.genreRepo.FindOne(ctx, []any{
		clause.Eq{Column: "slug", Value: slug},
	}, nil)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.ResultNotFound("Genre")
		}
		s.logger.Error("Failed to find genre for deletion", "error", err)
		return response.ResultErrDb(err)
	}

	if err := s.genreRepo.DeleteSoft(ctx, []any{
		clause.Eq{Column: "id", Value: genre.ID},
	}); err != nil {
		s.logger.Error("Failed to delete genre", "error", err)
		return response.ResultErrDb(err)
	}

	return response.ResultSuccess("Genre deleted successfully", nil)
}
