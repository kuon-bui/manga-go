package authorservice

import (
	"context"
	"errors"
	"manga-go/internal/app/api/common/response"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (s *AuthorService) GetAuthor(ctx context.Context, id uuid.UUID) response.Result {
	author, err := s.authorRepo.FindOne(ctx, []any{
		clause.Eq{Column: "id", Value: id},
	}, nil)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.ResultNotFound("Author")
		}
		s.logger.Error("Failed to find author", "error", err)
		return response.ResultErrDb(err)
	}

	return response.ResultSuccess("Author retrieved successfully", author)
}
