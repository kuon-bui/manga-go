package comicservice

import (
	"context"
	"manga-go/internal/app/api/common/response"
	"manga-go/internal/pkg/common"

	"github.com/google/uuid"
)

func (s *ComicService) ListFollowedComics(ctx context.Context, userID uuid.UUID, paging *common.Paging) response.Result {
	follows, total, err := s.comicFollowRepo.FindPaginatedByUserID(ctx, userID, paging)
	if err != nil {
		s.logger.Error("Failed to list followed comics", "error", err)
		return response.ResultErrDb(err)
	}

	return response.ResultPaginationData(follows, total, "Followed comics retrieved successfully")
}
