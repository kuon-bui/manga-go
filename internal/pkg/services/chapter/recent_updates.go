package chapterserivce

import (
	"context"
	"manga-go/internal/app/api/common/response"
	chapterrequest "manga-go/internal/pkg/request/chapter"
	chapterresponse "manga-go/internal/pkg/response/chapter"
)

func (s *ChapterService) GetRecentUpdates(ctx context.Context, req *chapterrequest.RecentUpdatesRequest) response.Result {
	chapters, total, err := s.chapterRepo.FindRecentUpdates(ctx, &req.Paging)
	if err != nil {
		s.logger.Error("Failed to fetch recent updates", "error", err)
		return response.ResultErrDb(err)
	}

	responses := make([]*chapterresponse.RecentUpdateResponse, 0, len(chapters))
	for _, chapter := range chapters {
		resp := &chapterresponse.RecentUpdateResponse{
			Chapter: chapter,
			Title:   chapter.Comic,
		}

		resp.Chapter.Comic = nil

		responses = append(responses, resp)
	}

	return response.ResultPaginationData(responses, total, "Recent updates retrieved successfully")
}
