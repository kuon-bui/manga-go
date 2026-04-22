package comicservice

import (
	"context"
	"manga-go/internal/app/api/common/response"
	"manga-go/internal/pkg/common"
	comicrequest "manga-go/internal/pkg/request/comic"
)

func (s *ComicService) ListComics(ctx context.Context, req *comicrequest.ListComicsRequest) response.Result {

	comics, total, err := s.comicRepo.FindPaginatedWithFilters(ctx, req, map[string]common.MoreKeyOption{
		"Artists":          {},
		"Authors":          {},
		"Genres":           {},
		"Tags":             {},
		"TranslationGroup": {},
		"UploadedBy":       {},
	})
	if err != nil {
		s.logger.Error("Failed to list comics", "error", err)
		return response.ResultErrDb(err)
	}

	return response.ResultPaginationData(comics, total, "Comics retrieved successfully")
}
