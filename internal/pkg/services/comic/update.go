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
		"description":        req.Description,
		"thumbnail":          req.Thumbnail,
		"banner":             req.Banner,
		"type":               req.Type,
		"status":             req.Status,
		"artist":             req.Artist,
		"published_year":     req.PublishedYear,
	}

	if req.IsActive != nil {
		updateData["is_active"] = *req.IsActive
	}
	if req.IsHot != nil {
		updateData["is_hot"] = *req.IsHot
	}
	if req.IsFeatured != nil {
		updateData["is_featured"] = *req.IsFeatured
	}

	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&model.Comic{}).Where("id = ?", comic.ID).Updates(updateData).Error; err != nil {
			return err
		}

		if req.AuthorIDs != nil {
			var authors []model.Author
			for _, aid := range req.AuthorIDs {
				authors = append(authors, model.Author{SqlModel: common.SqlModel{ID: aid}})
			}
			if err := tx.Model(comic).Association("Authors").Replace(authors); err != nil {
				return err
			}
		}

		if req.GenreIDs != nil {
			var genres []model.Genre
			for _, gid := range req.GenreIDs {
				genres = append(genres, model.Genre{SqlModel: common.SqlModel{ID: gid}})
			}
			if err := tx.Model(comic).Association("Genres").Replace(genres); err != nil {
				return err
			}
		}

		if req.TagIDs != nil {
			var tags []model.Tag
			for _, tid := range req.TagIDs {
				tags = append(tags, model.Tag{SqlModel: common.SqlModel{ID: tid}})
			}
			if err := tx.Model(comic).Association("Tags").Replace(tags); err != nil {
				return err
			}
		}

		return nil
	})
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
