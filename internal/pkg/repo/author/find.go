package authorrepo

import (
	"manga-go/internal/pkg/common"
	"manga-go/internal/pkg/model"

	"github.com/samber/lo"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (r *AuthorRepo) FindByNamesWithTx(tx *gorm.DB, names []string, moreKeys map[string]common.MoreKeyOption) ([]*model.Author, error) {
	return r.FindAllWithTx(
		tx,
		[]any{
			clause.IN{
				Column: clause.Column{
					Name:  "name",
					Table: clause.CurrentTable,
				},
				Values: lo.Map(names, func(e string, _ int) any { return e }),
			},
		},
		moreKeys,
	)
}
