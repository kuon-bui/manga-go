package tagservice

import (
	"context"
	"errors"
	"manga-go/internal/app/api/common/response"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (s *TagService) DeleteTag(ctx context.Context, slug string) response.Result {
	tag, err := s.tagRepo.FindOne(ctx, []any{
		clause.Eq{Column: "slug", Value: slug},
	}, nil)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.ResultNotFound("Tag")
		}
		s.logger.Error("Failed to find tag for deletion", "error", err)
		return response.ResultErrDb(err)
	}

	if err := s.tagRepo.DeleteSoft(ctx, []any{
		clause.Eq{Column: "id", Value: tag.ID},
	}); err != nil {
		s.logger.Error("Failed to delete tag", "error", err)
		return response.ResultErrDb(err)
	}

	return response.ResultSuccess("Tag deleted successfully", nil)
}
