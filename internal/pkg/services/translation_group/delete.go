package translationgroupservice

import (
	"context"
	"errors"
	"manga-go/internal/app/api/common/response"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (s *TranslationGroupService) DeleteTranslationGroup(ctx context.Context, slug string) response.Result {
	_, err := s.translationGroupRepo.FindOne(ctx, []any{
		clause.Eq{Column: "slug", Value: slug},
	}, nil)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.ResultNotFound("TranslationGroup")
		}
		s.logger.Error("Failed to find translation group for deletion", "error", err)
		return response.ResultErrDb(err)
	}

	if err := s.translationGroupRepo.DeleteSoft(ctx, []any{
		clause.Eq{Column: "slug", Value: slug},
	}); err != nil {
		s.logger.Error("Failed to delete translation group", "error", err)
		return response.ResultErrDb(err)
	}

	return response.ResultSuccess("Translation group deleted successfully", nil)
}
