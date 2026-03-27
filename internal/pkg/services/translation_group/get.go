package translationgroupservice

import (
	"context"
	"errors"
	"manga-go/internal/app/api/common/response"
	"manga-go/internal/pkg/common"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (s *TranslationGroupService) GetTranslationGroup(ctx context.Context, slug string) response.Result {
	group, err := s.translationGroupRepo.FindOne(ctx, []any{
		clause.Eq{Column: "slug", Value: slug},
	}, map[string]common.MoreKeyOption{
		"Owner":   {},
		"Members": {},
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.ResultNotFound("TranslationGroup")
		}
		s.logger.Error("Failed to find translation group", "error", err)
		return response.ResultErrDb(err)
	}

	return response.ResultSuccess("Translation group retrieved successfully", group)
}
