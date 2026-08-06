package comicservice

import (
	"context"
	"manga-go/internal/app/api/common/response"
	"manga-go/internal/pkg/authorization"
	"manga-go/internal/pkg/common"
	comicrequest "manga-go/internal/pkg/request/comic"

	"gorm.io/gorm"
)

func (s *ComicService) ListComics(ctx context.Context, req *comicrequest.ListComicsRequest) response.Result {
	scopes := []func(*gorm.DB) *gorm.DB{}
	if authorization.ViewerFromContext(ctx).User == nil {
		scopes = append(scopes, s.comicRepo.IsPublished(true))
	}

	comics, total, err := s.comicRepo.FindPaginatedWithFilters(ctx, req, map[string]common.MoreKeyOption{
		"Artists":          {},
		"Authors":          {},
		"Genres":           {},
		"Tags":             {},
		"TranslationGroup": {},
		"UploadedBy":       {},
	}, scopes...)
	if err != nil {
		s.logger.Error("Failed to list comics", "error", err)
		return response.ResultErrDb(err)
	}

	return response.ResultPaginationData(comics, total, "Comics retrieved successfully")
}
