package comicrepo

import (
	"manga-go/internal/pkg/common"
	"manga-go/internal/pkg/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (r *ComicRepo) UpdateComicWithTransaction(tx *gorm.DB, id uuid.UUID, data map[string]any, associations map[string]any) error {
	// Update comic
	if err := r.UpdateWithTransaction(tx, []any{
		clause.Eq{Column: "id", Value: id},
	}, data); err != nil {
		return err
	}

	// Update associations
	for assocName, assocData := range associations {
		if err := tx.Model(&model.Comic{SqlModel: common.SqlModel{ID: id}}).Association(assocName).Replace(assocData); err != nil {
			return err
		}
	}

	return nil
}
