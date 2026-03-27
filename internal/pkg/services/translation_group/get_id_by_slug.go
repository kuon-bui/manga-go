package translationgroupservice

import (
	"context"

	"github.com/google/uuid"
)

func (s *TranslationGroupService) GetTranslationGroupIDBySlug(ctx context.Context, slug string) (uuid.UUID, error) {
	return s.translationGroupRepo.GetIdBySlug(ctx, slug)
}
