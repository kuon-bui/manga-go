package authorizationrevision

import (
	"context"
	"fmt"
	"time"

	"manga-go/internal/pkg/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const globalScope = "global"

type Repo struct {
	db *gorm.DB
}

func NewRepo(db *gorm.DB) *Repo {
	return &Repo{db: db}
}

func UserScope(userID uuid.UUID) string {
	return "user:" + userID.String()
}

func (r *Repo) Current(ctx context.Context, userID uuid.UUID) (uint64, uint64, error) {
	if err := r.ensureScopes(r.db.WithContext(ctx), globalScope, UserScope(userID)); err != nil {
		return 0, 0, err
	}

	var revisions []model.AuthorizationCacheRevision
	if err := r.db.WithContext(ctx).
		Where("scope IN ?", []string{globalScope, UserScope(userID)}).
		Find(&revisions).Error; err != nil {
		return 0, 0, err
	}

	versions := make(map[string]uint64, len(revisions))
	for _, revision := range revisions {
		versions[revision.Scope] = revision.Version
	}
	global, globalOK := versions[globalScope]
	user, userOK := versions[UserScope(userID)]
	if !globalOK || !userOK {
		return 0, 0, fmt.Errorf("authorization revisions are incomplete")
	}
	return global, user, nil
}

func (r *Repo) BumpGlobalTx(tx *gorm.DB) (uint64, error) {
	return r.bumpTx(tx, globalScope)
}

func (r *Repo) BumpUserTx(tx *gorm.DB, userID uuid.UUID) (uint64, error) {
	return r.bumpTx(tx, UserScope(userID))
}

func (r *Repo) bumpTx(tx *gorm.DB, scope string) (uint64, error) {
	if err := r.ensureScopes(tx, scope); err != nil {
		return 0, err
	}
	if err := tx.Model(&model.AuthorizationCacheRevision{}).
		Where("scope = ?", scope).
		Updates(map[string]any{
			"version":    gorm.Expr("version + 1"),
			"updated_at": time.Now().UTC(),
		}).Error; err != nil {
		return 0, err
	}

	var revision model.AuthorizationCacheRevision
	if err := tx.Where("scope = ?", scope).First(&revision).Error; err != nil {
		return 0, err
	}
	return revision.Version, nil
}

func (r *Repo) ensureScopes(db *gorm.DB, scopes ...string) error {
	revisions := make([]model.AuthorizationCacheRevision, 0, len(scopes))
	now := time.Now().UTC()
	for _, scope := range scopes {
		revisions = append(revisions, model.AuthorizationCacheRevision{
			Scope:     scope,
			Version:   1,
			UpdatedAt: now,
		})
	}
	return db.Clauses(clause.OnConflict{DoNothing: true}).Create(&revisions).Error
}
