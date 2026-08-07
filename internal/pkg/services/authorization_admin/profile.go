package authorizationadmin

import (
	"context"
	"errors"
	"fmt"
	"sort"

	"manga-go/internal/pkg/authorization"
	"manga-go/internal/pkg/model"

	"github.com/google/uuid"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var authorizationProfileCacheFailures = promauto.NewCounterVec(
	prometheus.CounterOpts{
		Name: "authorization_profile_cache_failures_total",
		Help: "Number of authorization profile cache failures by operation.",
	},
	[]string{"operation"},
)

func (s *Service) GetProfile(ctx context.Context, userID uuid.UUID) (*AuthorizationProfile, error) {
	globalVersion, userVersion, err := s.revisions.Current(ctx, userID)
	if err != nil {
		return nil, err
	}

	key := profileCacheKey(userID, globalVersion, userVersion)
	profile := &AuthorizationProfile{}
	cacheReadable := true
	hit, err := s.cache.Get(ctx, key, profile)
	if err != nil {
		cacheReadable = false
		s.logCacheError("read", key, err)
	} else if hit {
		return profile, nil
	}

	profile, err = s.calculateProfile(ctx, userID, globalVersion, userVersion)
	if err != nil {
		return nil, err
	}
	if cacheReadable {
		if err := s.cache.Set(ctx, key, profile, profileTTL); err != nil {
			s.logCacheError("write", key, err)
		}
	}
	return profile, nil
}

func (s *Service) calculateProfile(
	ctx context.Context,
	userID uuid.UUID,
	globalVersion uint64,
	userVersion uint64,
) (*AuthorizationProfile, error) {
	roleIDs, err := s.policyManager.RolesForUser(userID.String(), authorization.OrgPlatform)
	if err != nil {
		return nil, err
	}

	roles, err := s.loadRoleSummaries(ctx, roleIDs)
	if err != nil {
		return nil, err
	}

	permissions := make([]string, 0)
	for _, definition := range authorization.Catalog() {
		allowed, err := s.isEffective(ctx, userID.String(), definition)
		if err != nil {
			return nil, err
		}
		if allowed {
			permissions = append(permissions, definition.Name)
		}
	}
	sort.Strings(permissions)

	return &AuthorizationProfile{
		UserID:      userID,
		Roles:       roles,
		Permissions: permissions,
		Version:     fmt.Sprintf("g%d:u%d", globalVersion, userVersion),
	}, nil
}

func (s *Service) loadRoleSummaries(ctx context.Context, rawIDs []string) ([]RoleSummary, error) {
	roles := make([]RoleSummary, 0, len(rawIDs))
	if len(rawIDs) == 0 {
		return roles, nil
	}

	ids := make([]uuid.UUID, 0, len(rawIDs))
	for _, rawID := range rawIDs {
		id, err := uuid.Parse(rawID)
		if err != nil {
			return nil, fmt.Errorf("invalid role id %q in policy: %w", rawID, err)
		}
		ids = append(ids, id)
	}

	var models []model.Role
	if err := s.roleRepo.DB.WithContext(ctx).Where("id IN ?", ids).Find(&models).Error; err != nil {
		return nil, err
	}
	for _, role := range models {
		roles = append(roles, RoleSummary{
			ID:          role.ID,
			Name:        role.Name,
			Description: role.Description,
		})
	}
	sort.Slice(roles, func(i, j int) bool {
		return roles[i].Name < roles[j].Name
	})
	return roles, nil
}

func (s *Service) isEffective(
	ctx context.Context,
	subject string,
	definition authorization.PermissionDefinition,
) (bool, error) {
	for _, action := range definition.Grants {
		err := s.authorizer.Enforce(ctx, authorization.Request{
			Subject: subject,
			Org:     authorization.OrgPlatform,
			Action:  action,
			Object:  definition.Object,
			Context: authorization.CtxAny,
		})
		if errors.Is(err, authorization.ErrForbidden) {
			return false, nil
		}
		if err != nil {
			return false, err
		}
	}
	return true, nil
}

func profileCacheKey(userID uuid.UUID, globalVersion uint64, userVersion uint64) string {
	return fmt.Sprintf(
		"authorization:profile:g%d:u%d:user:%s",
		globalVersion,
		userVersion,
		userID.String(),
	)
}

func (s *Service) logCacheError(operation string, key string, err error) {
	authorizationProfileCacheFailures.WithLabelValues(operation).Inc()
	if s.logger == nil {
		return
	}
	s.logger.Warnw("authorization profile cache failure", "operation", operation, "key", key, "error", err)
}
