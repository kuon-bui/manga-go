package comicservice

import (
	"context"
	"errors"
	"manga-go/internal/app/api/common/response"
	"manga-go/internal/pkg/common"
	"manga-go/internal/pkg/model"
	comicrequest "manga-go/internal/pkg/request/comic"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (s *ComicService) UpdateComic(ctx context.Context, slug string, req *comicrequest.UpdateComicRequest) response.Result {
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

	updateData := map[string]any{
		"title":              req.Title,
		"slug":               req.Slug,
		"alternative_titles": common.StringSlice(req.AlternativeTitles),
		"type":               req.Type,
		"status":             req.Status,
	}

	if req.Description != nil {
		updateData["description"] = req.Description
	}

	if req.Thumbnail != nil {
		updateData["thumbnail"] = req.Thumbnail
	}

	if req.Banner != nil {
		updateData["banner"] = req.Banner
	}

	if req.PublishedYear != nil {
		updateData["published_year"] = req.PublishedYear
	}

	if req.IsHot != nil {
		updateData["is_hot"] = *req.IsHot
	}

	if req.IsFeatured != nil {
		updateData["is_featured"] = *req.IsFeatured
	}

	if req.AgeRating != "" {
		updateData["age_rating"] = req.AgeRating
	}

	associations := make(map[string]any)

	var genres []*model.Genre
	if len(req.GenreSlugs) > 0 {
		foundGenres, err := s.genreRepo.FindBySlugs(ctx, req.GenreSlugs, nil)
		if err != nil {
			s.logger.Error("Failed to find genres by slug", "error", err)
			return response.ResultErrDb(err)
		}

		genresBySlug := make(map[string]*model.Genre, len(foundGenres))
		for _, genre := range foundGenres {
			genresBySlug[genre.Slug] = genre
		}

		for _, genreSlug := range req.GenreSlugs {
			genre, ok := genresBySlug[genreSlug]
			if !ok {
				return response.ResultError("Genre slug not found: " + genreSlug)
			}
			genres = append(genres, genre)
		}
		associations["Genres"] = genres
	}

	var tags []*model.Tag
	if len(req.TagSlugs) > 0 {
		foundTags, err := s.tagRepo.FindBySlugs(ctx, req.TagSlugs, nil)
		if err != nil {
			s.logger.Error("Failed to find tags by slug", "error", err)
			return response.ResultErrDb(err)
		}

		tagsBySlug := make(map[string]*model.Tag, len(foundTags))
		for _, tag := range foundTags {
			tagsBySlug[tag.Slug] = tag
		}

		for _, tagSlug := range req.TagSlugs {
			tag, ok := tagsBySlug[tagSlug]
			if !ok {
				return response.ResultError("Tag slug not found: " + tagSlug)
			}
			tags = append(tags, tag)
		}
		associations["Tags"] = tags
	}

	if req.AuthorIDs != nil {
		var authors []*model.Author
		for _, aid := range req.AuthorIDs {
			authors = append(authors, &model.Author{SqlModel: common.SqlModel{ID: aid}})
		}
		associations["Authors"] = authors
	}

	err = s.comicRepo.UpdateComicWithTransaction(ctx, comic.ID, updateData, associations)
	if err != nil {
		s.logger.Error("Failed to update comic", "error", err)
		return response.ResultErrDb(err)
	}

	updated, err := s.comicRepo.FindOne(ctx, []any{
		clause.Eq{Column: "id", Value: comic.ID},
	}, map[string]common.MoreKeyOption{
		"Authors": {},
		"Genres":  {},
		"Tags":    {},
	})
	if err != nil {
		s.logger.Error("Failed to reload comic", "error", err)
		return response.ResultErrDb(err)
	}

	return response.ResultSuccess("Comic updated successfully", updated)
}
