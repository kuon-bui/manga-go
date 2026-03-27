package tagservice

import (
	"context"

	"github.com/google/uuid"
)

func (s *TagService) GetTagIDBySlug(ctx context.Context, slug string) (uuid.UUID, error) {
	return s.tagRepo.GetIdBySlug(ctx, slug)
}
