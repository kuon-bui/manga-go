package translationgroupservice

import (
	"context"
	"manga-go/internal/app/api/common/response"
	"manga-go/internal/pkg/common"
)

func (s *TranslationGroupService) ListTranslationGroups(ctx context.Context, paging *common.Paging) response.Result {
	groups, total, err := s.translationGroupRepo.FindPaginated(ctx, nil, paging, nil)
	if err != nil {
		s.logger.Error("Failed to list translation groups", "error", err)
		return response.ResultErrDb(err)
	}

	return response.ResultSuccess("Translation groups retrieved successfully", response.ResponsePaginationData(groups, total))
}
