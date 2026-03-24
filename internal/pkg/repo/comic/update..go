package comicrepo

import (
	"context"
	"manga-go/internal/pkg/common"
	"manga-go/internal/pkg/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (r *ComicRepo) UpdateComicWithTransaction(ctx context.Context, id uuid.UUID, data map[string]any, associations map[string]any) error {
	return r.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
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
	})
}
