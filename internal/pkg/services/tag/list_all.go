package tagservice

import (
	"context"
	"manga-go/internal/app/api/common/response"
)

func (s *TagService) ListAllTags(ctx context.Context) response.Result {
	tags, err := s.tagRepo.FindAll(ctx, nil, nil)
	if err != nil {
		s.logger.Error("Failed to list all tags", "error", err)
		return response.ResultErrDb(err)
	}

	return response.ResultSuccess("Tags retrieved successfully", tags)
}
