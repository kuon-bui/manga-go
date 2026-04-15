package comicfollowrepo

import (
	"context"
	"manga-go/internal/pkg/common"
	"manga-go/internal/pkg/model"
	"manga-go/internal/pkg/repo/base"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ComicFollowRepo struct {
	*base.BaseRepository[model.ComicFollow]
}

func NewComicFollowRepo(db *gorm.DB) *ComicFollowRepo {
	return &ComicFollowRepo{
		BaseRepository: &base.BaseRepository[model.ComicFollow]{
			DB: db,
		},
	}
}

func (r *ComicFollowRepo) FindByUserAndComic(ctx context.Context, userID, comicID uuid.UUID) (*model.ComicFollow, error) {
	return r.FindOne(ctx, []any{
		clause.Eq{Column: "user_id", Value: userID},
		clause.Eq{Column: "comic_id", Value: comicID},
	}, nil)
}

func (r *ComicFollowRepo) FindByUserAndComicWithUnscoped(ctx context.Context, userID, comicID uuid.UUID) (*model.ComicFollow, error) {
	return r.FindOneWithUnscoped(ctx, []any{
		clause.Eq{Column: "user_id", Value: userID},
		clause.Eq{Column: "comic_id", Value: comicID},
	}, nil)
}

func (r *ComicFollowRepo) RestoreByUserAndComic(ctx context.Context, userID, comicID uuid.UUID) error {
	return r.DB.WithContext(ctx).
		Model(&model.ComicFollow{}).
		Unscoped().
		Where(clause.Eq{Column: "user_id", Value: userID}).
		Where(clause.Eq{Column: "comic_id", Value: comicID}).
		Update("deleted_at", nil).Error
}

func (r *ComicFollowRepo) FindPaginatedByUserID(ctx context.Context, userID uuid.UUID, paging *common.Paging) ([]*model.ComicFollow, int64, error) {
	var follows []*model.ComicFollow
	var total int64

	db := r.DB.WithContext(ctx).
		Model(&model.ComicFollow{}).
		Where(clause.Eq{Column: "user_id", Value: userID}).
		Preload("Comic").
		Order("created_at DESC")

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	db = db.Scopes(r.WithPaginate(paging))

	if err := db.Find(&follows).Error; err != nil {
		return nil, 0, err
	}

	return follows, total, nil
}
