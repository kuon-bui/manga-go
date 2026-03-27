package tagservice

import (
	"context"
	"manga-go/internal/app/api/common/response"
	"manga-go/internal/pkg/model"
	tagrequest "manga-go/internal/pkg/request/tag"
)

func (s *TagService) CreateTag(ctx context.Context, req *tagrequest.CreateTagRequest) response.Result {
	tag := model.Tag{
		Name: req.Name,
		Slug: req.Slug,
	}

	if err := s.tagRepo.Create(ctx, &tag); err != nil {
		s.logger.Error("Failed to create tag", "error", err)
		return response.ResultErrDb(err)
	}

	return response.ResultSuccess("Tag created successfully", tag)
}
