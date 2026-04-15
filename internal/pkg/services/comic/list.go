package comicservice

import (
	"context"
	"errors"
	"manga-go/internal/app/api/common/response"
	"manga-go/internal/pkg/common"
	comicrepo "manga-go/internal/pkg/repo/comic"
	comicrequest "manga-go/internal/pkg/request/comic"
	"strings"

	"gorm.io/gorm"
)

func (s *ComicService) ListComics(ctx context.Context, req *comicrequest.ListComicsRequest) response.Result {
	filters := comicrepo.ListComicFilters{
		GenreSlugs: req.GenreSlugs,
		TagSlugs:   req.TagSlugs,
		Search:     strings.TrimSpace(req.Search),
	}

	if req.TranslationGroupSlug != "" {
		translationGroupID, err := s.translationGroupRepo.GetIdBySlug(ctx, req.TranslationGroupSlug)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return response.ResultNotFound("TranslationGroup")
			}
			s.logger.Error("Failed to get translation group by slug", "error", err)
			return response.ResultErrDb(err)
		}

		filters.TranslationGroupID = &translationGroupID
	}

	comics, total, err := s.comicRepo.FindPaginatedWithFilters(ctx, filters, &req.Paging, map[string]common.MoreKeyOption{
		"Artists": {},
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
