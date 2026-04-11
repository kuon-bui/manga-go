package permissionservice_test

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"manga-go/internal/pkg/common"
	"manga-go/internal/pkg/logger"
	"manga-go/internal/pkg/model"
	permissionrequest "manga-go/internal/pkg/request/permission"
	permissionservice "manga-go/internal/pkg/services/permission"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// MockPermissionRepo implements permissionservice.PermissionRepository using testify/mock.
type MockPermissionRepo struct {
	mock.Mock
}

func (m *MockPermissionRepo) Create(ctx context.Context, permission *model.Permission) error {
	args := m.Called(ctx, permission)
	return args.Error(0)
}

func (m *MockPermissionRepo) FindOne(ctx context.Context, conditions []any, moreKeys map[string]common.MoreKeyOption) (*model.Permission, error) {
	args := m.Called(ctx, conditions, moreKeys)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Permission), args.Error(1)
}

func (m *MockPermissionRepo) Update(ctx context.Context, conditions []any, data map[string]any) error {
	args := m.Called(ctx, conditions, data)
	return args.Error(0)
}

func (m *MockPermissionRepo) DeleteSoft(ctx context.Context, conditions []any) error {
	args := m.Called(ctx, conditions)
	return args.Error(0)
}

func (m *MockPermissionRepo) FindPaginated(ctx context.Context, conditions []any, paging *common.Paging, moreKeys map[string]common.MoreKeyOption) ([]*model.Permission, int64, error) {
	args := m.Called(ctx, conditions, paging, moreKeys)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*model.Permission), args.Get(1).(int64), args.Error(2)
}

func (m *MockPermissionRepo) FindAll(ctx context.Context, conditions []any, moreKeys map[string]common.MoreKeyOption) ([]*model.Permission, error) {
	args := m.Called(ctx, conditions, moreKeys)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.Permission), args.Error(1)
}

func newTestService(repo *MockPermissionRepo) *permissionservice.PermissionService {
	return permissionservice.NewPermissionServiceWithRepo(logger.NewLogger(), repo)
}

func samplePermission() *model.Permission {
	now := time.Now()
	return &model.Permission{
		SqlModel: common.SqlModel{
			ID:        uuid.New(),
			CreatedAt: &now,
			UpdatedAt: &now,
		},
		Name: "read:comics",
	}
}

// ---------------------------------------------------------------------------
// CreatePermission
// ---------------------------------------------------------------------------

func TestCreatePermission_Success(t *testing.T) {
	repo := new(MockPermissionRepo)
	svc := newTestService(repo)

	req := &permissionrequest.CreatePermissionRequest{Name: "write:comics"}
	repo.On("Create", mock.Anything, mock.AnythingOfType("*model.Permission")).Return(nil)

	result := svc.CreatePermission(context.Background(), req)

	assert.True(t, result.Success)
	assert.Equal(t, http.StatusOK, result.HttpStatus)
	assert.Equal(t, "Permission created successfully", result.Message)
	repo.AssertExpectations(t)
}

func TestCreatePermission_DBError(t *testing.T) {
	repo := new(MockPermissionRepo)
	svc := newTestService(repo)

	req := &permissionrequest.CreatePermissionRequest{Name: "write:comics"}
	repo.On("Create", mock.Anything, mock.AnythingOfType("*model.Permission")).Return(errors.New("db error"))

	result := svc.CreatePermission(context.Background(), req)

	assert.False(t, result.Success)
	assert.Equal(t, http.StatusInternalServerError, result.HttpStatus)
	repo.AssertExpectations(t)
}

// ---------------------------------------------------------------------------
// ListPermissions
// ---------------------------------------------------------------------------

func TestListPermissions_Success(t *testing.T) {
	repo := new(MockPermissionRepo)
	svc := newTestService(repo)

	permissions := []*model.Permission{samplePermission(), samplePermission()}
	paging := &common.Paging{Page: 1, Limit: 20}
	repo.On("FindPaginated", mock.Anything, mock.Anything, paging, mock.Anything).Return(permissions, int64(2), nil)

	result := svc.ListPermissions(context.Background(), paging)

	assert.True(t, result.Success)
	assert.Equal(t, http.StatusOK, result.HttpStatus)
	repo.AssertExpectations(t)
}

func TestListPermissions_DBError(t *testing.T) {
	repo := new(MockPermissionRepo)
	svc := newTestService(repo)

	paging := &common.Paging{Page: 1, Limit: 20}
	repo.On("FindPaginated", mock.Anything, mock.Anything, paging, mock.Anything).Return(nil, int64(0), errors.New("db error"))

	result := svc.ListPermissions(context.Background(), paging)

	assert.False(t, result.Success)
	assert.Equal(t, http.StatusInternalServerError, result.HttpStatus)
	repo.AssertExpectations(t)
}

// ---------------------------------------------------------------------------
// ListAllPermissions
// ---------------------------------------------------------------------------

func TestListAllPermissions_Success(t *testing.T) {
	repo := new(MockPermissionRepo)
	svc := newTestService(repo)

	permissions := []*model.Permission{samplePermission()}
	repo.On("FindAll", mock.Anything, mock.Anything, mock.Anything).Return(permissions, nil)

	result := svc.ListAllPermissions(context.Background())

	assert.True(t, result.Success)
	assert.Equal(t, http.StatusOK, result.HttpStatus)
	repo.AssertExpectations(t)
}

func TestListAllPermissions_DBError(t *testing.T) {
	repo := new(MockPermissionRepo)
	svc := newTestService(repo)

	repo.On("FindAll", mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("db error"))

	result := svc.ListAllPermissions(context.Background())

	assert.False(t, result.Success)
	assert.Equal(t, http.StatusInternalServerError, result.HttpStatus)
	repo.AssertExpectations(t)
}

// ---------------------------------------------------------------------------
// UpdatePermission
// ---------------------------------------------------------------------------

func TestUpdatePermission_Success(t *testing.T) {
	repo := new(MockPermissionRepo)
	svc := newTestService(repo)

	perm := samplePermission()
	req := &permissionrequest.UpdatePermissionRequest{Name: "delete:comics"}

	repo.On("FindOne", mock.Anything, mock.Anything, mock.Anything).Return(perm, nil)
	repo.On("Update", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	result := svc.UpdatePermission(context.Background(), perm.ID, req)

	assert.True(t, result.Success)
	assert.Equal(t, http.StatusOK, result.HttpStatus)
	assert.Equal(t, "Permission updated successfully", result.Message)
	repo.AssertExpectations(t)
}

func TestUpdatePermission_NotFound(t *testing.T) {
	repo := new(MockPermissionRepo)
	svc := newTestService(repo)

	req := &permissionrequest.UpdatePermissionRequest{Name: "delete:comics"}
	repo.On("FindOne", mock.Anything, mock.Anything, mock.Anything).Return(nil, gorm.ErrRecordNotFound)

	result := svc.UpdatePermission(context.Background(), uuid.New(), req)

	assert.False(t, result.Success)
	assert.Equal(t, http.StatusBadRequest, result.HttpStatus)
	assert.Equal(t, "Permission not found", result.Message)
	repo.AssertExpectations(t)
}

func TestUpdatePermission_UpdateDBError(t *testing.T) {
	repo := new(MockPermissionRepo)
	svc := newTestService(repo)

	perm := samplePermission()
	req := &permissionrequest.UpdatePermissionRequest{Name: "delete:comics"}

	repo.On("FindOne", mock.Anything, mock.Anything, mock.Anything).Return(perm, nil)
	repo.On("Update", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("update failed"))

	result := svc.UpdatePermission(context.Background(), perm.ID, req)

	assert.False(t, result.Success)
	assert.Equal(t, http.StatusInternalServerError, result.HttpStatus)
	repo.AssertExpectations(t)
}

// ---------------------------------------------------------------------------
// DeletePermission
// ---------------------------------------------------------------------------

func TestDeletePermission_Success(t *testing.T) {
	repo := new(MockPermissionRepo)
	svc := newTestService(repo)

	perm := samplePermission()
	repo.On("FindOne", mock.Anything, mock.Anything, mock.Anything).Return(perm, nil)
	repo.On("DeleteSoft", mock.Anything, mock.Anything).Return(nil)

	result := svc.DeletePermission(context.Background(), perm.ID)

	assert.True(t, result.Success)
	assert.Equal(t, http.StatusOK, result.HttpStatus)
	assert.Equal(t, "Permission deleted successfully", result.Message)
	repo.AssertExpectations(t)
}

func TestDeletePermission_NotFound(t *testing.T) {
	repo := new(MockPermissionRepo)
	svc := newTestService(repo)

	repo.On("FindOne", mock.Anything, mock.Anything, mock.Anything).Return(nil, gorm.ErrRecordNotFound)

	result := svc.DeletePermission(context.Background(), uuid.New())

	assert.False(t, result.Success)
	assert.Equal(t, http.StatusBadRequest, result.HttpStatus)
	assert.Equal(t, "Permission not found", result.Message)
	repo.AssertExpectations(t)
}

func TestDeletePermission_DeleteDBError(t *testing.T) {
	repo := new(MockPermissionRepo)
	svc := newTestService(repo)

	perm := samplePermission()
	repo.On("FindOne", mock.Anything, mock.Anything, mock.Anything).Return(perm, nil)
	repo.On("DeleteSoft", mock.Anything, mock.Anything).Return(errors.New("delete failed"))

	result := svc.DeletePermission(context.Background(), perm.ID)

	assert.False(t, result.Success)
	assert.Equal(t, http.StatusInternalServerError, result.HttpStatus)
	repo.AssertExpectations(t)
}
