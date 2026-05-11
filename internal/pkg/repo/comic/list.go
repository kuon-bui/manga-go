package comicrepo

import (
	"context"
	"manga-go/internal/pkg/common"
	"manga-go/internal/pkg/model"
	comicrequest "manga-go/internal/pkg/request/comic"

	"gorm.io/gorm/clause"
)

var allowedSortFields = map[string]string{
	"lastChapterAt": "comics.last_chapter_at",
	"createdAt":     "comics.created_at",
	"rating":        "COALESCE(cs.avg_rating, 0)",
	"followCount":   "COALESCE(cs.follow_count, 0)",
}

const statsJoin = "LEFT JOIN comic_stats cs ON cs.comic_id = comics.id"

const statsSelect = `comics.*,
	COALESCE(cs.follow_count, 0) AS follow_count,
	COALESCE(cs.rating_count, 0) AS rating_count,
	COALESCE(cs.chapter_count, 0) AS chapter_count,
	cs.avg_rating AS avg_rating`

func buildComicSortOrder(sortBy, order string) clause.OrderByColumn {
	sortExpr, ok := allowedSortFields[sortBy]
	if !ok {
		sortExpr = "created_at"
	}
	orderMap := map[string]bool{
		"desc": true,
		"DESC": true,
	}

	return clause.OrderByColumn{
		Column: clause.Column{
			Name:  sortExpr,
			Table: clause.CurrentTable,
		},
		Desc: orderMap[order],
	}
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
		Joins(statsJoin).
		Select(statsSelect).
		Distinct().
		Scopes(r.WithPaginate(&filters.Paging)).
		Order(buildComicSortOrder(filters.SortBy, filters.Order))

	query = r.ApplyPreloadMoreKeys(query, moreKeys)

	if err := query.Find(&comics).Error; err != nil {
		return nil, 0, err
	}

	return comics, total, nil
}

func (r *ComicRepo) FindOneWithStats(ctx context.Context, conditions []any, moreKeys map[string]common.MoreKeyOption) (*model.Comic, error) {
	var comic model.Comic
	db := r.DB.WithContext(ctx).
		Model(&model.Comic{}).
		Joins(statsJoin).
		Select(statsSelect)

	db = r.ApplyWhereConditions(db, conditions)
	db = r.ApplyPreloadMoreKeys(db, moreKeys)

	if err := db.First(&comic).Error; err != nil {
		return nil, err
	}
	return &comic, nil
}
