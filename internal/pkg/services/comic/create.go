package comicservice

import (
	"context"
	"manga-go/internal/app/api/common/response"
	"manga-go/internal/pkg/common"
	"manga-go/internal/pkg/model"
	comicrequest "manga-go/internal/pkg/request/comic"
)

func (s *ComicService) CreateComic(ctx context.Context, req *comicrequest.CreateComicRequest) response.Result {
	comic := model.Comic{
		Title:             req.Title,
		Slug:              req.Slug,
		AlternativeTitles: common.StringSlice(req.AlternativeTitles),
		Description:       req.Description,
		Thumbnail:         req.Thumbnail,
		Banner:            req.Banner,
		Type:              req.Type,
		Status:            req.Status,
		Artist:            req.Artist,
		PublishedYear:     req.PublishedYear,
	}

	if req.IsActive != nil {
		comic.IsActive = *req.IsActive
	} else {
		comic.IsActive = true
	}
	if req.IsHot != nil {
		comic.IsHot = *req.IsHot
	}
	if req.IsFeatured != nil {
		comic.IsFeatured = *req.IsFeatured
	}

	if comic.Type == "" {
		comic.Type = "manga"
	}
	if comic.Status == "" {
		comic.Status = "ongoing"
	}

	// Build author associations
	for _, aid := range req.AuthorIDs {
		comic.Authors = append(comic.Authors, model.Author{SqlModel: common.SqlModel{ID: aid}})
	}
	// Build genre associations
	for _, gid := range req.GenreIDs {
		comic.Genres = append(comic.Genres, model.Genre{SqlModel: common.SqlModel{ID: gid}})
	}
	// Build tag associations
	for _, tid := range req.TagIDs {
		comic.Tags = append(comic.Tags, model.Tag{SqlModel: common.SqlModel{ID: tid}})
	}

	if err := s.db.WithContext(ctx).Create(&comic).Error; err != nil {
		s.logger.Error("Failed to create comic", "error", err)
		return response.ResultErrDb(err)
	}

	return response.ResultSuccess("Comic created successfully", comic)
}
