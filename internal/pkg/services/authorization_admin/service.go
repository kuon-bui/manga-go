package authorizationadmin

import (
	"context"
	"time"

	"manga-go/internal/pkg/authorization"
	"manga-go/internal/pkg/logger"
	authorizationaudit "manga-go/internal/pkg/repo/authorization_audit"
	rolerepo "manga-go/internal/pkg/repo/role"
	userrepo "manga-go/internal/pkg/repo/user"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const profileTTL = 10 * time.Minute

type ProfileCache interface {
	Get(ctx context.Context, key string, target *AuthorizationProfile) (bool, error)
	Set(ctx context.Context, key string, profile *AuthorizationProfile, ttl time.Duration) error
	Delete(ctx context.Context, key string) error
}

type RevisionStore interface {
	Current(ctx context.Context, userID uuid.UUID) (global uint64, user uint64, err error)
	CurrentGlobal(ctx context.Context) (uint64, error)
	CurrentMany(ctx context.Context, userIDs []uuid.UUID) (uint64, map[uuid.UUID]uint64, error)
	BumpGlobalTx(tx *gorm.DB) (uint64, error)
	BumpUserTx(tx *gorm.DB, userID uuid.UUID) (uint64, error)
}

type Service struct {
	logger        *logger.Logger
	roleRepo      *rolerepo.RoleRepo
	userRepo      *userrepo.UserRepository
	policyManager *authorization.PolicyManager
	authorizer    *authorization.Authorizer
	auditRepo     *authorizationaudit.Repo
	revisions     RevisionStore
	cache         ProfileCache
}
