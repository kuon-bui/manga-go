package comicservice

import (
	"context"
	"manga-go/internal/app/api/common/response"
	"manga-go/internal/pkg/common"
	"manga-go/internal/pkg/constant"
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
		Type:              constant.ComicTypeManga,
		Status:            constant.ComicStatusOngoing, // Default status to ongoing
		IsActive:          true,                        // Default to active
		PublishedYear:     req.PublishedYear,
	}

	if req.ArtistId != nil {
		comic.ArtistId = req.ArtistId
	}

	if req.Type != "" {
		comic.Type = req.Type
	}

	// Build author associations
	for _, aid := range req.AuthorIDs {
		comic.Authors = append(comic.Authors, &model.Author{SqlModel: common.SqlModel{ID: aid}})
	}

	if len(req.GenreSlugs) > 0 {
		genres, err := s.genreRepo.FindBySlugs(ctx, req.GenreSlugs, nil)
		if err != nil {
			s.logger.Error("Failed to find genres by slug", "error", err)
			return response.ResultErrDb(err)
		}

		genresBySlug := make(map[string]*model.Genre, len(genres))
		for _, genre := range genres {
			genresBySlug[genre.Slug] = genre
		}

		for _, slug := range req.GenreSlugs {
			genre, ok := genresBySlug[slug]
			if !ok {
				return response.ResultError("Genre slug not found: " + slug)
			}
			comic.Genres = append(comic.Genres, genre)
		}
	}

	if len(req.TagSlugs) > 0 {
		tags, err := s.tagRepo.FindBySlugs(ctx, req.TagSlugs, nil)
		if err != nil {
			s.logger.Error("Failed to find tags by slug", "error", err)
			return response.ResultErrDb(err)
		}

		tagsBySlug := make(map[string]*model.Tag, len(tags))
		for _, tag := range tags {
			tagsBySlug[tag.Slug] = tag
		}

		for _, slug := range req.TagSlugs {
			tag, ok := tagsBySlug[slug]
			if !ok {
				return response.ResultError("Tag slug not found: " + slug)
			}
			comic.Tags = append(comic.Tags, tag)
		}
	}

	if err := s.comicRepo.Create(ctx, &comic); err != nil {
		s.logger.Error("Failed to create comic", "error", err)
		return response.ResultErrDb(err)
	}

	return response.ResultSuccess("Comic created successfully", comic)
}
