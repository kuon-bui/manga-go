package comicrepo

import (
	"context"
	"fmt"
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
	Status             string // ongoing | completed | hiatus | cancelled
	SortBy             string // lastChapterAt | createdAt | rating | followCount
	Order              string // asc | desc
}

var allowedSortFields = map[string]string{
	"lastChapterAt": "comics.last_chapter_at",
	"createdAt":     "comics.created_at",
	"rating":        "(SELECT COALESCE(AVG(r.score), 0) FROM ratings r WHERE r.comic_id = comics.id AND r.deleted_at IS NULL)",
	"followCount":   "(SELECT COUNT(*) FROM comic_follows cf WHERE cf.comic_id = comics.id AND cf.deleted_at IS NULL)",
}

const statsSelect = `comics.*,
	(SELECT COUNT(*) FROM comic_follows cf WHERE cf.comic_id = comics.id AND cf.deleted_at IS NULL) AS follow_count,
	(SELECT COUNT(*) FROM ratings r WHERE r.comic_id = comics.id AND r.deleted_at IS NULL) AS rating_count,
	(SELECT COUNT(*) FROM chapters ch WHERE ch.comic_id = comics.id AND ch.deleted_at IS NULL) AS chapter_count,
	(SELECT AVG(r.score) FROM ratings r WHERE r.comic_id = comics.id AND r.deleted_at IS NULL) AS avg_rating`

func (r *ComicRepo) FindPaginatedWithFilters(ctx context.Context, filters ListComicFilters, paging *common.Paging, moreKeys map[string]common.MoreKeyOption) ([]*model.Comic, int64, error) {
	var comics []*model.Comic
	var total int64

	countQuery := r.buildFilteredQuery(ctx, filters)
	if err := countQuery.Distinct("comics.id").Count(&total).Error; err != nil {
		return nil, 0, err
	}

	query := r.buildFilteredQuery(ctx, filters).
		Select(statsSelect).
		Distinct()

	// Apply sorting
	if sortField, ok := allowedSortFields[filters.SortBy]; ok {
		order := "DESC"
		if filters.Order == "asc" {
			order = "ASC"
		}
		query = query.Order(fmt.Sprintf("%s %s NULLS LAST", sortField, order))
	} else {
		query = query.Order("comics.created_at DESC")
	}

	query = r.ApplyPreloadMoreKeys(query, moreKeys)
	query = query.Scopes(r.WithPaginate(paging))

	if err := query.Find(&comics).Error; err != nil {
		return nil, 0, err
	}

	return comics, total, nil
}

func (r *ComicRepo) FindOneWithStats(ctx context.Context, conditions []any, moreKeys map[string]common.MoreKeyOption) (*model.Comic, error) {
	var comic model.Comic
	db := r.DB.WithContext(ctx).
		Model(&model.Comic{}).
		Select(statsSelect)

	db = r.ApplyWhereConditions(db, conditions)
	db = r.ApplyPreloadMoreKeys(db, moreKeys)

	if err := db.First(&comic).Error; err != nil {
		return nil, err
	}
	return &comic, nil
}

func (r *ComicRepo) buildFilteredQuery(ctx context.Context, filters ListComicFilters) *gorm.DB {
	db := r.DB.WithContext(ctx).Model(&model.Comic{})

	if filters.TranslationGroupID != nil {
		db = db.Where("comics.translation_group_id = ?", *filters.TranslationGroupID)
	}

	if filters.Status != "" {
		db = db.Where("comics.status = ?", filters.Status)
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
