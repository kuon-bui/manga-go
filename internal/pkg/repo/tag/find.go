package tagrepo

import (
	"context"
	"manga-go/internal/pkg/common"
	"manga-go/internal/pkg/model"

	"github.com/samber/lo"
	"gorm.io/gorm/clause"
)

func (r *TagRepo) FindBySlugs(ctx context.Context, slugs []string, moreKeys map[string]common.MoreKeyOption) ([]*model.Tag, error) {
	return r.FindAll(
		ctx,
		[]any{
			clause.IN{
				Column: clause.Column{
					Name:  "slug",
					Table: clause.CurrentTable,
				},
				Values: lo.Map(slugs, func(e string, _ int) any { return e }),
			},
		},
		moreKeys,
	)
}
