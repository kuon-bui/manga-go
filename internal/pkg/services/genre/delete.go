package genreservice

import (
	"context"
	"errors"
	"manga-go/internal/app/api/common/response"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (s *GenreService) DeleteGenre(ctx context.Context, id uuid.UUID) response.Result {
	_, err := s.genreRepo.FindOne(ctx, []any{
		clause.Eq{Column: "id", Value: id},
	}, nil)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.ResultNotFound("Genre")
		}
		s.logger.Error("Failed to find genre for deletion", "error", err)
		return response.ResultErrDb(err)
	}

	if err := s.genreRepo.DeleteSoft(ctx, []any{
		clause.Eq{Column: "id", Value: id},
	}); err != nil {
		s.logger.Error("Failed to delete genre", "error", err)
		return response.ResultErrDb(err)
	}

	return response.ResultSuccess("Genre deleted successfully", nil)
}
