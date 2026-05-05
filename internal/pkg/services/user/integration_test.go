//go:build integration

package userservice

import (
	"context"
	"testing"
	"time"

	"manga-go/internal/pkg/logger"
	rolerepo "manga-go/internal/pkg/repo/role"
	userrepo "manga-go/internal/pkg/repo/user"
	"manga-go/internal/pkg/testutil"

	"github.com/google/uuid"
)

func newUserServiceIntegration(t *testing.T) *UserService {
	t.Helper()

	db := testutil.NewSQLiteDB(t)
	tx := db.Begin()
	if tx.Error != nil {
		t.Fatalf("failed to begin transaction: %v", tx.Error)
	}
	t.Cleanup(func() {
		_ = tx.Rollback().Error
	})

	testutil.MustSyncSchemas(t, tx,
		&testutil.User{},
		&testutil.Role{},
		&testutil.UserRole{},
	)

	return &UserService{
		logger:   logger.NewLogger(),
		userRepo: userrepo.NewUserRepository(tx, nil),
		roleRepo: rolerepo.NewRoleRepo(tx),
	}
}

func TestUserServiceIntegrationAssignAndRemoveRole(t *testing.T) {
	s := newUserServiceIntegration(t)
	ctx := context.Background()
	now := time.Now()

	userID := uuid.New()
	roleOneID := uuid.New()
	roleTwoID := uuid.New()

	if err := s.userRepo.DB.Exec(
		"INSERT INTO users (id, name, email, password, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)",
		userID,
		"integration-user",
		"integration@example.com",
		"hashed-password",
		now,
		now,
	).Error; err != nil {
		t.Fatalf("failed to seed user: %v", err)
	}

	if err := s.userRepo.DB.Exec(
		"INSERT INTO roles (id, name, created_at, updated_at) VALUES (?, ?, ?, ?), (?, ?, ?, ?)",
		roleOneID,
		"editor",
		now,
		now,
		roleTwoID,
		"translator",
		now,
		now,
	).Error; err != nil {
		t.Fatalf("failed to seed roles: %v", err)
	}

	assignRes := s.AssignRoles(ctx, userID, []uuid.UUID{roleOneID, roleTwoID})
	if !assignRes.Success {
		t.Fatalf("expected assign roles success, got: %s", assignRes.Message)
	}

	var assignedCount int64
	if err := s.userRepo.DB.Table("users_roles").Where("user_id = ?", userID).Count(&assignedCount).Error; err != nil {
		t.Fatalf("failed to count assigned user roles: %v", err)
	}
	if assignedCount != 2 {
		t.Fatalf("expected 2 assigned roles, got %d", assignedCount)
	}

	removeRes := s.RemoveRole(ctx, userID, roleOneID)
	if !removeRes.Success {
		t.Fatalf("expected remove role success, got: %s", removeRes.Message)
	}

	var remainingCount int64
	if err := s.userRepo.DB.Table("users_roles").Where("user_id = ?", userID).Count(&remainingCount).Error; err != nil {
		t.Fatalf("failed to count remaining user roles: %v", err)
	}
	if remainingCount != 1 {
		t.Fatalf("expected 1 remaining role, got %d", remainingCount)
	}
}
