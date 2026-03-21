package genreservice

import (
	"context"
	"errors"
	"manga-go/internal/app/api/common/response"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (s *GenreService) GetGenre(ctx context.Context, id uuid.UUID) response.Result {
	genre, err := s.genreRepo.FindOne(ctx, []any{
		clause.Eq{Column: "id", Value: id},
	}, nil)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.ResultNotFound("Genre")
		}
		s.logger.Error("Failed to find genre", "error", err)
		return response.ResultErrDb(err)
	}

	return response.ResultSuccess("Genre retrieved successfully", genre)
}
