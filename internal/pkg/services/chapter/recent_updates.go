package chapterserivce

import (
	"context"
	"manga-go/internal/app/api/common/response"
	chapterrequest "manga-go/internal/pkg/request/chapter"
	"time"

	"github.com/google/uuid"
)

type RecentUpdateTitle struct {
	ID         uuid.UUID `json:"id"`
	Name       string    `json:"name"`
	CoverImage *string   `json:"coverImage"`
}

type RecentUpdateChapter struct {
	ID        uuid.UUID  `json:"id"`
	Number    string     `json:"number"`
	Name      string     `json:"name"`
	CreatedAt *time.Time `json:"createdAt"`
}

type RecentUpdateResponse struct {
	Title   RecentUpdateTitle   `json:"title"`
	Chapter RecentUpdateChapter `json:"chapter"`
}

func (s *ChapterService) GetRecentUpdates(ctx context.Context, req *chapterrequest.RecentUpdatesRequest) response.Result {
	chapters, total, err := s.chapterRepo.FindRecentUpdates(ctx, &req.Paging)
	if err != nil {
		s.logger.Error("Failed to fetch recent updates", "error", err)
		return response.ResultErrDb(err)
	}

	responses := make([]*RecentUpdateResponse, 0, len(chapters))
	for _, chapter := range chapters {
		resp := &RecentUpdateResponse{
			Chapter: RecentUpdateChapter{
				ID:        chapter.ID,
				Number:    chapter.Number,
				Name:      chapter.Title,
				CreatedAt: chapter.CreatedAt,
			},
		}

		if chapter.Comic != nil {
			resp.Title = RecentUpdateTitle{
				ID:         chapter.Comic.ID,
				Name:       chapter.Comic.Title,
				CoverImage: chapter.Comic.Thumbnail,
			}
		}

		responses = append(responses, resp)
	}

	return response.ResultPaginationData(responses, total, "Recent updates retrieved successfully")
}
