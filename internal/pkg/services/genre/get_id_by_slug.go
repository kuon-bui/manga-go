package genreservice

import (
	"context"

	"github.com/google/uuid"
)

func (s *GenreService) GetGenreIDBySlug(ctx context.Context, slug string) (uuid.UUID, error) {
	return s.genreRepo.GetIdBySlug(ctx, slug)
}
