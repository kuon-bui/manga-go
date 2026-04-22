package comicrepo

import (
	"context"
	"fmt"
	"manga-go/internal/pkg/common"
	"manga-go/internal/pkg/model"
	comicrequest "manga-go/internal/pkg/request/comic"
)

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

func buildComicSortOrder(sortBy, order string) string {
	sortExpr, ok := allowedSortFields[sortBy]
	if !ok {
		return "comics.created_at DESC"
	}

	return fmt.Sprintf("%s %s NULLS LAST", sortExpr, order)
}

func (r *ComicRepo) FindPaginatedWithFilters(ctx context.Context, filters *comicrequest.ListComicsRequest, moreKeys map[string]common.MoreKeyOption) ([]*model.Comic, int64, error) {
	var comics []*model.Comic
	var total int64

	countQuery := r.DB.WithContext(ctx).
		Model(&model.Comic{}).
		Scopes(applyComicFilters(filters))
	if err := countQuery.Distinct("comics.id").Count(&total).Error; err != nil {
		return nil, 0, err
	}

	query := r.DB.WithContext(ctx).
		Model(&model.Comic{}).
		Scopes(applyComicFilters(filters)).
		Select(statsSelect).
		Distinct()

	query = query.Order(buildComicSortOrder(filters.SortBy, filters.Order))

	query = r.ApplyPreloadMoreKeys(query, moreKeys)
	query = query.Scopes(r.WithPaginate(&filters.Paging))

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
