package chapterserivce

import (
	"context"

	"github.com/google/uuid"
)

func (s *ChapterService) GetChapterIDBySlug(ctx context.Context, slug string) (uuid.UUID, error) {
	return s.chapterRepo.GetIdBySlug(ctx, slug)
}
