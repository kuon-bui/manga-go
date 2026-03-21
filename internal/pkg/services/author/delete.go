package authorservice

import (
	"context"
	"errors"
	"manga-go/internal/app/api/common/response"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (s *AuthorService) DeleteAuthor(ctx context.Context, id uuid.UUID) response.Result {
	_, err := s.authorRepo.FindOne(ctx, []any{
		clause.Eq{Column: "id", Value: id},
	}, nil)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.ResultNotFound("Author")
		}
		s.logger.Error("Failed to find author for deletion", "error", err)
		return response.ResultErrDb(err)
	}

	if err := s.authorRepo.DeleteSoft(ctx, []any{
		clause.Eq{Column: "id", Value: id},
	}); err != nil {
		s.logger.Error("Failed to delete author", "error", err)
		return response.ResultErrDb(err)
	}

	return response.ResultSuccess("Author deleted successfully", nil)
}
