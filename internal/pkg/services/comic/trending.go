package comicservice

import (
	"context"
	"manga-go/internal/app/api/common/response"
	"manga-go/internal/pkg/model"

	"github.com/google/uuid"
)

type TrendingChapterInfo struct {
	Name   string    `json:"name"`
	Number string    `json:"number"`
	ID     uuid.UUID `json:"id"`
}

type TrendingComicResponse struct {
	ID            uuid.UUID            `json:"id"`
	Title         string               `json:"title"`
	CoverImage    *string              `json:"coverImage"`
	Author        string               `json:"author"`
	Synopsis      *string              `json:"synopsis"`
	Genres        []*model.Genre       `json:"genres"`
	LatestChapter *TrendingChapterInfo `json:"latestChapter"`
	Views         int                  `json:"views"`
}

func (s *ComicService) GetTrendingComics(ctx context.Context, limit int) response.Result {
	if limit <= 0 {
		limit = 5
	}
	if limit > 50 {
		limit = 50
	}

	comics, err := s.comicRepo.FindTrending(ctx, limit)
	if err != nil {
		s.logger.Error("Failed to fetch trending comics", "error", err)
		return response.ResultErrDb(err)
	}

	responses := make([]*TrendingComicResponse, 0, len(comics))
	for _, comic := range comics {
		resp := &TrendingComicResponse{
			ID:         comic.ID,
			Title:      comic.Title,
			CoverImage: comic.Thumbnail,
			Synopsis:   comic.Description,
			Genres:     comic.Genres,
			Views:      comic.FollowCount,
		}

		if len(comic.Authors) > 0 && comic.Authors[0] != nil {
			resp.Author = comic.Authors[0].Name
		}

		if comic.LatestChapter != nil {
			resp.LatestChapter = &TrendingChapterInfo{
				Name:   comic.LatestChapter.Title,
				Number: comic.LatestChapter.Number,
				ID:     comic.LatestChapter.ID,
			}
		}

		responses = append(responses, resp)
	}

	return response.ResultSuccess("Trending comics retrieved successfully", responses)
}
