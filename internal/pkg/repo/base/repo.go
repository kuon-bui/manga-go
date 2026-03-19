package base

import (
	"context"
	"manga-go/internal/pkg/common"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type BaseRepository[T ModelInterface] struct {
	DB *gorm.DB
}

type ModelInterface interface {
	TableName() string
}

func (r *BaseRepository[T]) ApplyPreloadMoreKeys(db *gorm.DB, moreKeys map[string]common.MoreKeyOption) *gorm.DB {
	for relation, option := range moreKeys {
		db = db.Preload(relation, func(tx *gorm.DB) *gorm.DB {
			if option.Custom != nil {
				return option.Custom(tx)
			}

			if option.Unscoped {
				tx = tx.Unscoped()
			}
			if option.Order != nil {
				tx = tx.Order(*option.Order)
			}
			if option.Where != nil {
				tx = tx.Where(*option.Where)
			}
			if len(option.Select) > 0 {
				tx = tx.Select(option.Select)
			}
			if option.Limit != nil {
				tx = tx.Limit(*option.Limit)
			}
			return tx
		})
	}
	return db
}

func (r *BaseRepository[T]) FindOne(ctx context.Context, conditions []any, moreKeys map[string]common.MoreKeyOption) (*T, error) {
	var model T
	db := r.DB.WithContext(ctx)

	// Xử lý JOIN và WHERE
	for _, condition := range conditions {
		switch c := condition.(type) {
		case common.JoinExpr:
			db = db.Joins(c.SQL, c.Vars...)
		default:
			db = db.Where(condition)
		}
	}

	// PRELOAD relationships with options
	db = r.ApplyPreloadMoreKeys(db, moreKeys)

	// Lấy bản ghi đầu tiên
	if err := db.First(&model).Error; err != nil {
		return nil, err
	}
	return &model, nil
}

func (r *BaseRepository[T]) FindOneWithUnscoped(ctx context.Context, conditions []any, moreKeys map[string]common.MoreKeyOption) (*T, error) {
	var model T
	db := r.DB.WithContext(ctx)

	// WHERE conditions
	for _, condition := range conditions {
		db = db.Where(condition)
	}

	// PRELOAD relationships with options
	db = r.ApplyPreloadMoreKeys(db, moreKeys)

	// Lấy bản ghi đầu tiên
	if err := db.Unscoped().First(&model).Error; err != nil {
		return nil, err
	}
	return &model, nil
}

func (r *BaseRepository[T]) FindOneWithTransaction(tx *gorm.DB, conditions []any, moreKeys map[string]common.MoreKeyOption) (*T, error) {
	var model T
	db := tx

	// Xử lý JOIN và WHERE
	for _, condition := range conditions {
		switch c := condition.(type) {
		case common.JoinExpr:
			db = db.Joins(c.SQL, c.Vars...)
		default:
			db = db.Where(condition)
		}
	}

	// PRELOAD relationships with options
	db = r.ApplyPreloadMoreKeys(db, moreKeys)

	// Lấy bản ghi đầu tiên
	if err := db.First(&model).Error; err != nil {
		return nil, err
	}
	return &model, nil
}

func (r *BaseRepository[T]) FindAll(ctx context.Context, conditions []any, moreKeys map[string]common.MoreKeyOption) ([]*T, error) {
	var models []*T
	db := r.DB.WithContext(ctx)

	// where
	for _, condition := range conditions {
		db = db.Where(condition)
	}

	// PRELOAD relationships with options
	db = r.ApplyPreloadMoreKeys(db, moreKeys)

	// Lấy tất cả bản ghi
	if err := db.Find(&models).Error; err != nil {
		return nil, err
	}

	return models, nil
}

func (r *BaseRepository[T]) FindAllWithUnscoped(ctx context.Context, conditions []any, moreKeys map[string]common.MoreKeyOption) ([]*T, error) {
	var models []*T
	db := r.DB.WithContext(ctx)

	// where
	for _, condition := range conditions {
		db = db.Where(condition)
	}

	// PRELOAD relationships with options
	db = r.ApplyPreloadMoreKeys(db, moreKeys)

	// Lấy tất cả bản ghi
	if err := db.Unscoped().Find(&models).Error; err != nil {
		return nil, err
	}

	return models, nil
}

func (r *BaseRepository[T]) Create(ctx context.Context, t *T) error {
	return r.DB.WithContext(ctx).Create(t).Error
}

func (r *BaseRepository[T]) CreateWithTransaction(tx *gorm.DB, t *T) error {
	return tx.Create(t).Error
}

func (r *BaseRepository[T]) CreateList(ctx context.Context, t []*T) error {
	if len(t) == 0 {
		return nil
	}

	return r.DB.WithContext(ctx).Create(t).Error
}

func (r *BaseRepository[T]) CreateListWithTransaction(tx *gorm.DB, t []*T) error {
	if len(t) == 0 {
		return nil
	}

	return tx.Create(t).Error
}

func (r *BaseRepository[T]) Update(ctx context.Context, conditions []any, data map[string]any) error {

	db := r.DB.WithContext(ctx).Model(new(T))

	// where
	for _, condition := range conditions {
		db = db.Where(condition)
	}

	return db.Updates(data).Error
}

func (r *BaseRepository[T]) UpdateWithTransaction(tx *gorm.DB, conditions []any, data map[string]any) error {

	db := tx.Model(new(T))

	// where
	for _, condition := range conditions {
		db = db.Where(condition)
	}

	return db.Updates(data).Error
}

func (r *BaseRepository[T]) UpdateUnscopedWithTransaction(tx *gorm.DB, conditions []any, data map[string]any) error {

	db := tx.Model(new(T))

	// where
	for _, condition := range conditions {
		db = db.Where(condition)
	}

	return db.Unscoped().Updates(data).Error
}

func (r *BaseRepository[T]) DeletePermanently(ctx context.Context, conditions []any) error {
	t := new(T)
	db := r.DB.WithContext(ctx).Model(t)

	// where
	for _, condition := range conditions {
		db = db.Where(condition)
	}

	return db.Unscoped().Delete(t).Error
}

func (r *BaseRepository[T]) CountAll(ctx context.Context, conditions []any) (int64, error) {
	var count int64
	t := new(T)
	db := r.DB.WithContext(ctx).Model(t)

	for _, condition := range conditions {
		db = db.Where(condition)
	}

	if err := db.Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *BaseRepository[T]) DeletePermanentlyWithTransaction(tx *gorm.DB, conditions []any) error {
	t := new(T)
	db := tx.Model(t)

	// where
	for _, condition := range conditions {
		db = db.Where(condition)
	}

	return db.Unscoped().Delete(t).Error
}

func (r *BaseRepository[T]) DeleteSoft(ctx context.Context, conditions []any) error {
	t := new(T)
	db := r.DB.WithContext(ctx).Model(t)

	// where
	for _, condition := range conditions {
		db = db.Where(condition)
	}

	return db.Delete(t).Error
}

func (r *BaseRepository[T]) DeleteSoftWithTransaction(tx *gorm.DB, conditions []any) error {
	t := new(T)
	db := tx.Model(t)

	// where
	for _, condition := range conditions {
		db = db.Where(condition)
	}

	return db.Delete(t).Error
}

func (r *BaseRepository[T]) Upsert(ctx context.Context, entity *T, conflictColumns []string, updateColumns []string) error {
	return r.DB.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   toClauseColumns(conflictColumns),
		DoUpdates: clause.AssignmentColumns(updateColumns),
	}).Create(entity).Error
}

func (r *BaseRepository[T]) UpsertWithTransaction(tx *gorm.DB, entity *T, conflictColumns []string, updateColumns []string) error {
	return tx.Clauses(clause.OnConflict{
		Columns:   toClauseColumns(conflictColumns),
		DoUpdates: clause.AssignmentColumns(updateColumns),
	}).Create(entity).Error
}

func (r *BaseRepository[T]) UpsertMany(ctx context.Context, entities []*T, conflictColumns []string, updateColumns []string) error {
	return r.DB.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   toClauseColumns(conflictColumns),
		DoUpdates: clause.AssignmentColumns(updateColumns),
	}).Create(entities).Error
}

func (r *BaseRepository[T]) UpsertManyWithTransaction(tx *gorm.DB, entities []*T, conflictColumns []string, updateColumns []string) error {
	return tx.Clauses(clause.OnConflict{
		Columns:   toClauseColumns(conflictColumns),
		DoUpdates: clause.AssignmentColumns(updateColumns),
	}).Create(entities).Error
}

func toClauseColumns(cols []string) []clause.Column {
	result := make([]clause.Column, len(cols))
	for i, c := range cols {
		result[i] = clause.Column{Name: c}
	}
	return result
}
