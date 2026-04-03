package comicservice

import (
	"context"
	"manga-go/internal/app/api/common/response"
	"manga-go/internal/pkg/common"
)

func (s *ComicService) ListComics(ctx context.Context, paging *common.Paging) response.Result {
	comics, total, err := s.comicRepo.FindPaginated(ctx, nil, paging, map[string]common.MoreKeyOption{
		"Authors": {},
		"Genres":  {},
		"Tags":    {},
	})
	if err != nil {
		s.logger.Error("Failed to list comics", "error", err)
		return response.ResultErrDb(err)
	}

	return response.ResultPaginationData(comics, total, "Comics retrieved successfully")
}
