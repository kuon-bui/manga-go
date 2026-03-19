package base

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (r *BaseRepository[T]) NotSoftDelete(db *gorm.DB) *gorm.DB {
	var tableName string
	if t, ok := any(new(T)).(ModelInterface); ok {
		tableName = t.TableName()
	} else {
		panic("T does not implement TableName() string")
	}

	return db.Where(tableName + ".deleted_at IS NULL")
}

func (r *BaseRepository[T]) LoadAllAssociations(db *gorm.DB) *gorm.DB {
	return db.Preload(clause.Associations)
}

func (r *BaseRepository[T]) IsActive(isActive bool) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("active = ?", isActive)
	}
}
