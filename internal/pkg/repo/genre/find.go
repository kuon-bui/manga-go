package genrerepo

import (
	"context"
	"manga-go/internal/pkg/common"
	"manga-go/internal/pkg/model"

	"github.com/samber/lo"
	"gorm.io/gorm/clause"
)

func (r *GenreRepo) FindBySlugs(ctx context.Context, slugs []string, moreKey map[string]common.MoreKeyOption) ([]*model.Genre, error) {
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
		moreKey,
	)
}
