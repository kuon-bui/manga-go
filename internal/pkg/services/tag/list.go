package tagservice

import (
	"context"
	"manga-go/internal/app/api/common/response"
	"manga-go/internal/pkg/common"
)

func (s *TagService) ListTags(ctx context.Context, paging *common.Paging) response.Result {
	tags, total, err := s.tagRepo.FindPaginated(ctx, nil, paging, nil)
	if err != nil {
		s.logger.Error("Failed to list tags", "error", err)
		return response.ResultErrDb(err)
	}

	return response.ResultPaginationData(tags, total, "Tags retrieved successfully")
}
