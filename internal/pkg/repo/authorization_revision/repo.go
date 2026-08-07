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
	global, users, err := r.CurrentMany(ctx, []uuid.UUID{userID})
	if err != nil {
		return 0, 0, err
	}
	return global, users[userID], nil
}

func (r *Repo) CurrentGlobal(ctx context.Context) (uint64, error) {
	if err := r.ensureScopes(r.db.WithContext(ctx), globalScope); err != nil {
		return 0, err
	}
	var revision model.AuthorizationCacheRevision
	if err := r.db.WithContext(ctx).Where("scope = ?", globalScope).First(&revision).Error; err != nil {
		return 0, err
	}
	return revision.Version, nil
}

func (r *Repo) CurrentMany(ctx context.Context, userIDs []uuid.UUID) (uint64, map[uuid.UUID]uint64, error) {
	scopes := make([]string, 0, len(userIDs)+1)
	scopes = append(scopes, globalScope)
	for _, userID := range userIDs {
		scopes = append(scopes, UserScope(userID))
	}
	if err := r.ensureScopes(r.db.WithContext(ctx), scopes...); err != nil {
		return 0, nil, err
	}

	var revisions []model.AuthorizationCacheRevision
	if err := r.db.WithContext(ctx).
		Where("scope IN ?", scopes).
		Find(&revisions).Error; err != nil {
		return 0, nil, err
	}

	versions := make(map[string]uint64, len(revisions))
	for _, revision := range revisions {
		versions[revision.Scope] = revision.Version
	}
	global, globalOK := versions[globalScope]
	if !globalOK {
		return 0, nil, fmt.Errorf("global authorization revision is missing")
	}
	users := make(map[uuid.UUID]uint64, len(userIDs))
	for _, userID := range userIDs {
		version, ok := versions[UserScope(userID)]
		if !ok {
			return 0, nil, fmt.Errorf("authorization revision is missing for user %s", userID)
		}
		users[userID] = version
	}
	return global, users, nil
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
