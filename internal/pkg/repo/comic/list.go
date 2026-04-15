package comicrepo

import (
	"context"
	"manga-go/internal/pkg/common"
	"manga-go/internal/pkg/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ListComicFilters struct {
	TranslationGroupID *uuid.UUID
	GenreSlugs         []string
	TagSlugs           []string
	Search             string
}

func (r *ComicRepo) FindPaginatedWithFilters(ctx context.Context, filters ListComicFilters, paging *common.Paging, moreKeys map[string]common.MoreKeyOption) ([]*model.Comic, int64, error) {
	var comics []*model.Comic
	var total int64

	countQuery := r.buildFilteredQuery(ctx, filters)
	if err := countQuery.Distinct("comics.id").Count(&total).Error; err != nil {
		return nil, 0, err
	}

	query := r.buildFilteredQuery(ctx, filters).
		Select("comics.*").
		Distinct()

	query = r.ApplyPreloadMoreKeys(query, moreKeys)
	query = query.Scopes(r.WithPaginate(paging))

	if err := query.Find(&comics).Error; err != nil {
		return nil, 0, err
	}

	return comics, total, nil
}

func (r *ComicRepo) buildFilteredQuery(ctx context.Context, filters ListComicFilters) *gorm.DB {
	db := r.DB.WithContext(ctx).Model(&model.Comic{})

	if filters.TranslationGroupID != nil {
		db = db.Where("comics.translation_group_id = ?", *filters.TranslationGroupID)
	}

	if filters.Search != "" {
		searchPattern := "%" + filters.Search + "%"
		db = db.Where(
			`(
				comics.title ILIKE ?
				OR EXISTS (
					SELECT 1
					FROM jsonb_array_elements_text(COALESCE(comics.alternative_titles, '[]'::jsonb)) AS alt(title)
					WHERE alt.title ILIKE ?
				)
			)`,
			searchPattern,
			searchPattern,
		)
	}

	if len(filters.GenreSlugs) > 0 {
		db = db.
			Joins("JOIN comic_genres ON comic_genres.comic_id = comics.id").
			Joins("JOIN genres ON genres.id = comic_genres.genre_id").
			Where("genres.slug IN ?", filters.GenreSlugs)
	}

	if len(filters.TagSlugs) > 0 {
		db = db.
			Joins("JOIN comic_tags ON comic_tags.comic_id = comics.id").
			Joins("JOIN tags ON tags.id = comic_tags.tag_id").
			Where("tags.slug IN ?", filters.TagSlugs)
	}

	return db
}
