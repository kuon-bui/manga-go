package base

import (
	"manga-go/internal/pkg/common"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (r *BaseRepository[T]) NotSoftDelete(db *gorm.DB) *gorm.DB {
	return db.Where(clause.Eq{
		Column: clause.Column{Name: "deleted_at", Table: clause.CurrentTable},
		Value:  nil,
	})
}

func (r *BaseRepository[T]) LoadAllAssociations(db *gorm.DB) *gorm.DB {
	return db.Preload(clause.Associations)
}

func (r *BaseRepository[T]) IsPublished(isPublished bool) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("is_published = ?", isPublished)
	}
}

func (r *BaseRepository[T]) WithPaginate(paging *common.Paging) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if paging == nil {
			return db
		}

		if paging.GetLimit() <= 0 {
			return db
		}

		return db.Offset(paging.GetOffset()).Limit(paging.GetLimit())
	}
}
