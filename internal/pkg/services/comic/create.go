package comicservice

import (
	"context"
	"manga-go/internal/app/api/common/response"
	"manga-go/internal/pkg/common"
	"manga-go/internal/pkg/constant"
	"manga-go/internal/pkg/model"
	comicrequest "manga-go/internal/pkg/request/comic"
	"manga-go/internal/pkg/utils"
	"runtime/debug"
)

func (s *ComicService) CreateComic(ctx context.Context, req *comicrequest.CreateComicRequest) (result response.Result) {
	tx := s.gormDb.Begin().WithContext(ctx)
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			common.ShowDebugTrace("create comic", debug.Stack())
			result = response.ResultErrInternal(r.(error))
		} else if result.IsError() {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	comic := model.Comic{
		Title:             req.Title,
		Slug:              req.Slug,
		AlternativeTitles: common.StringSlice(req.AlternativeTitles),
		Description:       req.Description,
		Thumbnail:         req.Thumbnail,
		Banner:            req.Banner,
		Type:              constant.ComicTypeManga,
		Status:            constant.ComicStatusOngoing, // Default status to ongoing
		AgeRating:         constant.AgeRatingAll,       // Default age rating to all
		IsPublished:       false,                       // Default to unpublished
		PublishedYear:     req.PublishedYear,
	}

	if req.Type != "" {
		comic.Type = req.Type
	}

	if req.AgeRating != "" {
		comic.AgeRating = req.AgeRating
	}

	currentUser, err := utils.GetCurrentUserFormContext(ctx)
	if err != nil {
		s.logger.Error("Failed to get current user from context", "error", err)
		return response.ResultErrInternal(err)
	}

	comic.UploadedByID = &currentUser.ID

	if currentUser.TranslationGroupID != nil {
		comic.TranslationGroupID = currentUser.TranslationGroupID
	}

	if req.AuthorNames != nil {
		authors, err := s.resolveOrCreateAuthorsByNames(tx, req.AuthorNames)
		if err != nil {
			s.logger.Error("Failed to resolve authors by name", "error", err)
			return response.ResultErrDb(err)
		}
		comic.Authors = authors
	}

	if req.ArtistNames != nil {
		artists, err := s.resolveOrCreateAuthorsByNames(tx, req.ArtistNames)
		if err != nil {
			s.logger.Error("Failed to resolve artists by name", "error", err)
			return response.ResultErrDb(err)
		}
		comic.Artists = artists
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
