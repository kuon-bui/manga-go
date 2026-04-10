package comicservice

import (
	"context"
	"errors"
	"manga-go/internal/app/api/common/response"
	comicrequest "manga-go/internal/pkg/request/comic"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (s *ComicService) UpdateComicStatus(ctx context.Context, slug string, req *comicrequest.UpdateComicStatusRequest) response.Result {
	comic, err := s.comicRepo.FindOne(ctx, []any{
		clause.Eq{Column: "slug", Value: slug},
	}, nil)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.ResultNotFound("Comic")
		}
		s.logger.Error("Failed to find comic", "error", err)
		return response.ResultErrDb(err)
	}

	if err := s.comicRepo.Update(ctx, []any{
		clause.Eq{Column: "id", Value: comic.ID},
	}, map[string]any{
		"status": req.Status,
	}); err != nil {
		s.logger.Error("Failed to update comic status", "error", err)
		return response.ResultErrDb(err)
	}

	comic.Status = req.Status
	return response.ResultSuccess("Comic status updated successfully", comic)
}
