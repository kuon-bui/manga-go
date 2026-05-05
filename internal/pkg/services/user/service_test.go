package userservice

import (
	"context"
	"net/http"
	"testing"
	"time"

	"manga-go/internal/pkg/config"
	"manga-go/internal/pkg/hash"
	"manga-go/internal/pkg/logger"
	rolerepo "manga-go/internal/pkg/repo/role"
	userrepo "manga-go/internal/pkg/repo/user"
	userrequest "manga-go/internal/pkg/request/user"
	"manga-go/internal/pkg/testutil"

	"github.com/glebarez/sqlite"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func newUserService(t *testing.T, createTables bool) *UserService {
	t.Helper()

	db, err := gorm.Open(sqlite.Open(testutil.NewSQLiteDSN()), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open sqlite db: %v", err)
	}

	if createTables {
		testutil.MustSyncSchemas(t, db,
			&testutil.User{},
			&testutil.Role{},
			&testutil.UserRole{},
		)
	}

	return &UserService{
		logger:   logger.NewLogger(),
		userRepo: userrepo.NewUserRepository(db, nil),
		roleRepo: rolerepo.NewRoleRepo(db),
		config: &config.Config{
			ResetPassword: config.ResetPasswordConfig{
				TokenExpiryMinutes: 30,
				ResetPasswordURL:   "https://example.com/reset?token=%s",
			},
		},
	}
}

func seedUser(t *testing.T, s *UserService, email string, password string) uuid.UUID {
	t.Helper()

	id := uuid.New()
	err := s.userRepo.DB.Exec(
		"INSERT INTO users (id, name, email, password, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)",
		id.String(),
		"Test User",
		email,
		hash.HashPassword(password),
		time.Now(),
		time.Now(),
	).Error
	if err != nil {
		t.Fatalf("failed to seed user: %v", err)
	}

	return id
}

func TestCreateAccountReturnsErrorWhenEmailExists(t *testing.T) {
	t.Parallel()

	s := newUserService(t, true)
	seedUser(t, s, "john@example.com", "secret123")

	res := s.CreateAccount(context.Background(), &userrequest.CreateUserRequest{
		Name:     "John",
		Email:    "john@example.com",
		Password: "secret123",
	})

	if res.Success {
		t.Fatalf("expected failure result")
	}
	if res.HttpStatus != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", res.HttpStatus)
	}
	if res.Message != "User with this email already exists" {
		t.Fatalf("unexpected message: %s", res.Message)
	}
}

func TestCreateAccountReturnsDbErrorWhenTableMissing(t *testing.T) {
	t.Parallel()

	s := newUserService(t, false)
	res := s.CreateAccount(context.Background(), &userrequest.CreateUserRequest{
		Name:     "John",
		Email:    "john@example.com",
		Password: "secret123",
	})

	if res.Success {
		t.Fatalf("expected failure result")
	}
	if res.HttpStatus != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", res.HttpStatus)
	}
	if res.Message != "database error" {
		t.Fatalf("unexpected message: %s", res.Message)
	}
	if res.Error == nil {
		t.Fatalf("expected non-nil error")
	}
}

func TestSignInReturnsInvalidCredentialsOnWrongPassword(t *testing.T) {
	t.Parallel()

	s := newUserService(t, true)
	seedUser(t, s, "john@example.com", "correct-password")

	_, _, res := s.SignIn(context.Background(), &userrequest.SignInRequest{
		Email:    "john@example.com",
		Password: "wrong-password",
	})

	if res.Success {
		t.Fatalf("expected failure result")
	}
	if res.HttpStatus != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", res.HttpStatus)
	}
	if res.Message != "Invalid email or password" {
		t.Fatalf("unexpected message: %s", res.Message)
	}
}

func TestSignInReturnsDbErrorWhenTableMissing(t *testing.T) {
	t.Parallel()

	s := newUserService(t, false)
	_, _, res := s.SignIn(context.Background(), &userrequest.SignInRequest{
		Email:    "john@example.com",
		Password: "any-password",
	})

	if res.Success {
		t.Fatalf("expected failure result")
	}
	if res.HttpStatus != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", res.HttpStatus)
	}
	if res.Message != "database error" {
		t.Fatalf("unexpected message: %s", res.Message)
	}
	if res.Error == nil {
		t.Fatalf("expected non-nil error")
	}
}

func TestGetUserRolesReturnsNotFoundWhenUserMissing(t *testing.T) {
	t.Parallel()

	s := newUserService(t, true)
	res := s.GetUserRoles(context.Background(), uuid.New())

	if res.Success {
		t.Fatalf("expected failure result")
	}
	if res.HttpStatus != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", res.HttpStatus)
	}
	if res.Message != "User not found" {
		t.Fatalf("unexpected message: %s", res.Message)
	}
}

func TestGetUserRolesReturnsSuccessForExistingUser(t *testing.T) {
	t.Parallel()

	s := newUserService(t, true)
	userID := seedUser(t, s, "roles@example.com", "secret123")

	res := s.GetUserRoles(context.Background(), userID)

	if !res.Success {
		t.Fatalf("expected success result")
	}
	if res.HttpStatus != http.StatusOK {
		t.Fatalf("expected status 200, got %d", res.HttpStatus)
	}
	if res.Message != "User roles retrieved successfully" {
		t.Fatalf("unexpected message: %s", res.Message)
	}
}

func TestAssignRolesReturnsNotFoundWhenUserMissing(t *testing.T) {
	t.Parallel()

	s := newUserService(t, true)
	res := s.AssignRoles(context.Background(), uuid.New(), []uuid.UUID{uuid.New()})

	if res.Success {
		t.Fatalf("expected failure result")
	}
	if res.HttpStatus != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", res.HttpStatus)
	}
	if res.Message != "User not found" {
		t.Fatalf("unexpected message: %s", res.Message)
	}
}

func TestAssignRolesReturnsRoleNotFoundWhenRoleMissing(t *testing.T) {
	t.Parallel()

	s := newUserService(t, true)
	userID := seedUser(t, s, "alice@example.com", "secret123")

	res := s.AssignRoles(context.Background(), userID, []uuid.UUID{uuid.New()})

	if res.Success {
		t.Fatalf("expected failure result")
	}
	if res.HttpStatus != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", res.HttpStatus)
	}
	if res.Message != "Role not found" {
		t.Fatalf("unexpected message: %s", res.Message)
	}
}

func TestRemoveRoleReturnsNotFoundWhenUserMissing(t *testing.T) {
	t.Parallel()

	s := newUserService(t, true)
	res := s.RemoveRole(context.Background(), uuid.New(), uuid.New())

	if res.Success {
		t.Fatalf("expected failure result")
	}
	if res.HttpStatus != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", res.HttpStatus)
	}
	if res.Message != "User not found" {
		t.Fatalf("unexpected message: %s", res.Message)
	}
}

func TestRemoveRoleReturnsNotFoundWhenRoleMissing(t *testing.T) {
	t.Parallel()

	s := newUserService(t, true)
	userID := seedUser(t, s, "remove-role@example.com", "secret123")

	res := s.RemoveRole(context.Background(), userID, uuid.New())

	if res.Success {
		t.Fatalf("expected failure result")
	}
	if res.HttpStatus != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", res.HttpStatus)
	}
	if res.Message != "Role not found" {
		t.Fatalf("unexpected message: %s", res.Message)
	}
}

func TestRequestResetPasswordReturnsInvalidRequestWhenEmailEmpty(t *testing.T) {
	t.Parallel()

	s := newUserService(t, true)
	res := s.RequestResetPassword(context.Background(), "")

	if res.Success {
		t.Fatalf("expected failure result")
	}
	if res.HttpStatus != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", res.HttpStatus)
	}
	if res.Message != "invalid request" {
		t.Fatalf("unexpected message: %s", res.Message)
	}
	if res.Error == nil {
		t.Fatalf("expected non-nil error")
	}
}

func TestRequestResetPasswordReturnsDbErrorWhenUserMissing(t *testing.T) {
	t.Parallel()

	s := newUserService(t, true)
	res := s.RequestResetPassword(context.Background(), "unknown@example.com")

	if res.Success {
		t.Fatalf("expected failure result")
	}
	if res.HttpStatus != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", res.HttpStatus)
	}
	if res.Message != "database error" {
		t.Fatalf("unexpected message: %s", res.Message)
	}
	if res.Error == nil {
		t.Fatalf("expected non-nil error")
	}
}

func TestRequestResetPasswordReturnsDbErrorWhenTableMissing(t *testing.T) {
	t.Parallel()

	s := newUserService(t, false)
	res := s.RequestResetPassword(context.Background(), "john@example.com")

	if res.Success {
		t.Fatalf("expected failure result")
	}
	if res.HttpStatus != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", res.HttpStatus)
	}
	if res.Message != "database error" {
		t.Fatalf("unexpected message: %s", res.Message)
	}
	if res.Error == nil {
		t.Fatalf("expected non-nil error")
	}
}

func TestResetPasswordReturnsInvalidRequestWhenTokenInvalid(t *testing.T) {
	t.Parallel()

	s := newUserService(t, true)
	res := s.ResetPassword(context.Background(), userrequest.ResetPasswordRequest{
		Token:       "invalid-token",
		NewPassword: "new-password-123",
	})

	if res.Success {
		t.Fatalf("expected failure result")
	}
	if res.HttpStatus != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", res.HttpStatus)
	}
	if res.Message != "invalid request" {
		t.Fatalf("unexpected message: %s", res.Message)
	}
	if res.Error == nil {
		t.Fatalf("expected non-nil error")
	}
}

func TestResetPasswordReturnsDbErrorWhenTableMissing(t *testing.T) {
	t.Parallel()

	s := newUserService(t, false)
	res := s.ResetPassword(context.Background(), userrequest.ResetPasswordRequest{
		Token:       "token",
		NewPassword: "new-password-123",
	})

	if res.Success {
		t.Fatalf("expected failure result")
	}
	if res.HttpStatus != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", res.HttpStatus)
	}
	if res.Message != "database error" {
		t.Fatalf("unexpected message: %s", res.Message)
	}
	if res.Error == nil {
		t.Fatalf("expected non-nil error")
	}
}
