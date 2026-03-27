package comicservice

import (
	"context"

	"github.com/google/uuid"
)

func (s *ComicService) GetComicIDBySlug(ctx context.Context, slug string) (uuid.UUID, error) {
	return s.comicRepo.GetIdBySlug(ctx, slug)
}
