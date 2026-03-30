package comicservice

import (
	"context"

	"github.com/google/uuid"
)

func (s *ComicService) GetComicIDBySlug(ctx context.Context, slug string) (uuid.UUID, error) {
	return s.comicRepo.GetIdBySlug(ctx, slug)
}

// GetComicIDAndGroupIDBySlug returns the comic ID and its translation group ID for the given slug.
func (s *ComicService) GetComicIDAndGroupIDBySlug(ctx context.Context, slug string) (uuid.UUID, *uuid.UUID, error) {
	return s.comicRepo.GetIdAndGroupIdBySlug(ctx, slug)
}
