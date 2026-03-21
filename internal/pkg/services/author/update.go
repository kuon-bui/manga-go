package authorservice

import (
	"context"
	"errors"
	"manga-go/internal/app/api/common/response"
	authorrequest "manga-go/internal/pkg/request/author"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (s *AuthorService) UpdateAuthor(ctx context.Context, id uuid.UUID, req *authorrequest.UpdateAuthorRequest) response.Result {
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

	if err := s.authorRepo.Update(ctx, []any{
		clause.Eq{Column: "id", Value: author.ID},
	}, map[string]any{
		"name": req.Name,
	}); err != nil {
		s.logger.Error("Failed to update author", "error", err)
		return response.ResultErrDb(err)
	}

	author.Name = req.Name
	return response.ResultSuccess("Author updated successfully", author)
}
