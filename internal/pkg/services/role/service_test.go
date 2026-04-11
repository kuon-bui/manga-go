package roleservice_test

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"manga-go/internal/pkg/common"
	"manga-go/internal/pkg/logger"
	"manga-go/internal/pkg/model"
	rolerequest "manga-go/internal/pkg/request/role"
	roleservice "manga-go/internal/pkg/services/role"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// MockRoleRepo implements roleservice.RoleRepository using testify/mock.
type MockRoleRepo struct {
	mock.Mock
}

func (m *MockRoleRepo) Create(ctx context.Context, role *model.Role) error {
	args := m.Called(ctx, role)
	return args.Error(0)
}

func (m *MockRoleRepo) FindOne(ctx context.Context, conditions []any, moreKeys map[string]common.MoreKeyOption) (*model.Role, error) {
	args := m.Called(ctx, conditions, moreKeys)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Role), args.Error(1)
}

func (m *MockRoleRepo) Update(ctx context.Context, conditions []any, data map[string]any) error {
	args := m.Called(ctx, conditions, data)
	return args.Error(0)
}

func (m *MockRoleRepo) DeleteSoft(ctx context.Context, conditions []any) error {
	args := m.Called(ctx, conditions)
	return args.Error(0)
}

func (m *MockRoleRepo) FindPaginated(ctx context.Context, conditions []any, paging *common.Paging, moreKeys map[string]common.MoreKeyOption) ([]*model.Role, int64, error) {
	args := m.Called(ctx, conditions, paging, moreKeys)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*model.Role), args.Get(1).(int64), args.Error(2)
}

func (m *MockRoleRepo) FindAll(ctx context.Context, conditions []any, moreKeys map[string]common.MoreKeyOption) ([]*model.Role, error) {
	args := m.Called(ctx, conditions, moreKeys)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.Role), args.Error(1)
}

func (m *MockRoleRepo) AssignPermissions(ctx context.Context, roleID uuid.UUID, perms []*model.Permission) error {
	args := m.Called(ctx, roleID, perms)
	return args.Error(0)
}

func (m *MockRoleRepo) RemovePermission(ctx context.Context, roleID uuid.UUID, perm *model.Permission) error {
	args := m.Called(ctx, roleID, perm)
	return args.Error(0)
}

// MockPermissionRepo implements roleservice.PermissionRepository using testify/mock.
type MockPermissionRepo struct {
	mock.Mock
}

func (m *MockPermissionRepo) FindOne(ctx context.Context, conditions []any, moreKeys map[string]common.MoreKeyOption) (*model.Permission, error) {
	args := m.Called(ctx, conditions, moreKeys)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Permission), args.Error(1)
}

func newTestService(roleRepo *MockRoleRepo, permRepo *MockPermissionRepo) *roleservice.RoleService {
	return roleservice.NewRoleServiceWithRepos(logger.NewLogger(), roleRepo, permRepo)
}

func sampleRole() *model.Role {
	now := time.Now()
	return &model.Role{
		SqlModel: common.SqlModel{
			ID:        uuid.New(),
			CreatedAt: &now,
			UpdatedAt: &now,
		},
		Name: "admin",
	}
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
// CreateRole
// ---------------------------------------------------------------------------

func TestCreateRole_Success(t *testing.T) {
	roleRepo := new(MockRoleRepo)
	svc := newTestService(roleRepo, new(MockPermissionRepo))

	req := &rolerequest.CreateRoleRequest{Name: "editor"}
	roleRepo.On("Create", mock.Anything, mock.AnythingOfType("*model.Role")).Return(nil)

	result := svc.CreateRole(context.Background(), req)

	assert.True(t, result.Success)
	assert.Equal(t, http.StatusOK, result.HttpStatus)
	assert.Equal(t, "Role created successfully", result.Message)
	roleRepo.AssertExpectations(t)
}

func TestCreateRole_DBError(t *testing.T) {
	roleRepo := new(MockRoleRepo)
	svc := newTestService(roleRepo, new(MockPermissionRepo))

	req := &rolerequest.CreateRoleRequest{Name: "editor"}
	roleRepo.On("Create", mock.Anything, mock.AnythingOfType("*model.Role")).Return(errors.New("db error"))

	result := svc.CreateRole(context.Background(), req)

	assert.False(t, result.Success)
	assert.Equal(t, http.StatusInternalServerError, result.HttpStatus)
	roleRepo.AssertExpectations(t)
}

// ---------------------------------------------------------------------------
// GetRole
// ---------------------------------------------------------------------------

func TestGetRole_Success(t *testing.T) {
	roleRepo := new(MockRoleRepo)
	svc := newTestService(roleRepo, new(MockPermissionRepo))

	role := sampleRole()
	roleRepo.On("FindOne", mock.Anything, mock.Anything, mock.Anything).Return(role, nil)

	result := svc.GetRole(context.Background(), role.ID)

	assert.True(t, result.Success)
	assert.Equal(t, http.StatusOK, result.HttpStatus)
	assert.Equal(t, "Role retrieved successfully", result.Message)
	roleRepo.AssertExpectations(t)
}

func TestGetRole_NotFound(t *testing.T) {
	roleRepo := new(MockRoleRepo)
	svc := newTestService(roleRepo, new(MockPermissionRepo))

	roleRepo.On("FindOne", mock.Anything, mock.Anything, mock.Anything).Return(nil, gorm.ErrRecordNotFound)

	result := svc.GetRole(context.Background(), uuid.New())

	assert.False(t, result.Success)
	assert.Equal(t, http.StatusBadRequest, result.HttpStatus)
	assert.Equal(t, "Role not found", result.Message)
	roleRepo.AssertExpectations(t)
}

func TestGetRole_DBError(t *testing.T) {
	roleRepo := new(MockRoleRepo)
	svc := newTestService(roleRepo, new(MockPermissionRepo))

	roleRepo.On("FindOne", mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("timeout"))

	result := svc.GetRole(context.Background(), uuid.New())

	assert.False(t, result.Success)
	assert.Equal(t, http.StatusInternalServerError, result.HttpStatus)
	roleRepo.AssertExpectations(t)
}

// ---------------------------------------------------------------------------
// ListRoles
// ---------------------------------------------------------------------------

func TestListRoles_Success(t *testing.T) {
	roleRepo := new(MockRoleRepo)
	svc := newTestService(roleRepo, new(MockPermissionRepo))

	roles := []*model.Role{sampleRole(), sampleRole()}
	paging := &common.Paging{Page: 1, Limit: 20}
	roleRepo.On("FindPaginated", mock.Anything, mock.Anything, paging, mock.Anything).Return(roles, int64(2), nil)

	result := svc.ListRoles(context.Background(), paging)

	assert.True(t, result.Success)
	assert.Equal(t, http.StatusOK, result.HttpStatus)
	roleRepo.AssertExpectations(t)
}

func TestListRoles_DBError(t *testing.T) {
	roleRepo := new(MockRoleRepo)
	svc := newTestService(roleRepo, new(MockPermissionRepo))

	paging := &common.Paging{Page: 1, Limit: 20}
	roleRepo.On("FindPaginated", mock.Anything, mock.Anything, paging, mock.Anything).Return(nil, int64(0), errors.New("db error"))

	result := svc.ListRoles(context.Background(), paging)

	assert.False(t, result.Success)
	assert.Equal(t, http.StatusInternalServerError, result.HttpStatus)
	roleRepo.AssertExpectations(t)
}

// ---------------------------------------------------------------------------
// ListAllRoles
// ---------------------------------------------------------------------------

func TestListAllRoles_Success(t *testing.T) {
	roleRepo := new(MockRoleRepo)
	svc := newTestService(roleRepo, new(MockPermissionRepo))

	roles := []*model.Role{sampleRole()}
	roleRepo.On("FindAll", mock.Anything, mock.Anything, mock.Anything).Return(roles, nil)

	result := svc.ListAllRoles(context.Background())

	assert.True(t, result.Success)
	assert.Equal(t, http.StatusOK, result.HttpStatus)
	roleRepo.AssertExpectations(t)
}

func TestListAllRoles_DBError(t *testing.T) {
	roleRepo := new(MockRoleRepo)
	svc := newTestService(roleRepo, new(MockPermissionRepo))

	roleRepo.On("FindAll", mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("db error"))

	result := svc.ListAllRoles(context.Background())

	assert.False(t, result.Success)
	assert.Equal(t, http.StatusInternalServerError, result.HttpStatus)
	roleRepo.AssertExpectations(t)
}

// ---------------------------------------------------------------------------
// UpdateRole
// ---------------------------------------------------------------------------

func TestUpdateRole_Success(t *testing.T) {
	roleRepo := new(MockRoleRepo)
	svc := newTestService(roleRepo, new(MockPermissionRepo))

	role := sampleRole()
	req := &rolerequest.UpdateRoleRequest{Name: "superadmin"}

	roleRepo.On("FindOne", mock.Anything, mock.Anything, mock.Anything).Return(role, nil)
	roleRepo.On("Update", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	result := svc.UpdateRole(context.Background(), role.ID, req)

	assert.True(t, result.Success)
	assert.Equal(t, http.StatusOK, result.HttpStatus)
	assert.Equal(t, "Role updated successfully", result.Message)
	roleRepo.AssertExpectations(t)
}

func TestUpdateRole_NotFound(t *testing.T) {
	roleRepo := new(MockRoleRepo)
	svc := newTestService(roleRepo, new(MockPermissionRepo))

	req := &rolerequest.UpdateRoleRequest{Name: "superadmin"}
	roleRepo.On("FindOne", mock.Anything, mock.Anything, mock.Anything).Return(nil, gorm.ErrRecordNotFound)

	result := svc.UpdateRole(context.Background(), uuid.New(), req)

	assert.False(t, result.Success)
	assert.Equal(t, http.StatusBadRequest, result.HttpStatus)
	assert.Equal(t, "Role not found", result.Message)
	roleRepo.AssertExpectations(t)
}

func TestUpdateRole_UpdateDBError(t *testing.T) {
	roleRepo := new(MockRoleRepo)
	svc := newTestService(roleRepo, new(MockPermissionRepo))

	role := sampleRole()
	req := &rolerequest.UpdateRoleRequest{Name: "superadmin"}

	roleRepo.On("FindOne", mock.Anything, mock.Anything, mock.Anything).Return(role, nil)
	roleRepo.On("Update", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("update failed"))

	result := svc.UpdateRole(context.Background(), role.ID, req)

	assert.False(t, result.Success)
	assert.Equal(t, http.StatusInternalServerError, result.HttpStatus)
	roleRepo.AssertExpectations(t)
}

// ---------------------------------------------------------------------------
// DeleteRole
// ---------------------------------------------------------------------------

func TestDeleteRole_Success(t *testing.T) {
	roleRepo := new(MockRoleRepo)
	svc := newTestService(roleRepo, new(MockPermissionRepo))

	role := sampleRole()
	roleRepo.On("FindOne", mock.Anything, mock.Anything, mock.Anything).Return(role, nil)
	roleRepo.On("DeleteSoft", mock.Anything, mock.Anything).Return(nil)

	result := svc.DeleteRole(context.Background(), role.ID)

	assert.True(t, result.Success)
	assert.Equal(t, http.StatusOK, result.HttpStatus)
	assert.Equal(t, "Role deleted successfully", result.Message)
	roleRepo.AssertExpectations(t)
}

func TestDeleteRole_NotFound(t *testing.T) {
	roleRepo := new(MockRoleRepo)
	svc := newTestService(roleRepo, new(MockPermissionRepo))

	roleRepo.On("FindOne", mock.Anything, mock.Anything, mock.Anything).Return(nil, gorm.ErrRecordNotFound)

	result := svc.DeleteRole(context.Background(), uuid.New())

	assert.False(t, result.Success)
	assert.Equal(t, http.StatusBadRequest, result.HttpStatus)
	roleRepo.AssertExpectations(t)
}

func TestDeleteRole_DeleteDBError(t *testing.T) {
	roleRepo := new(MockRoleRepo)
	svc := newTestService(roleRepo, new(MockPermissionRepo))

	role := sampleRole()
	roleRepo.On("FindOne", mock.Anything, mock.Anything, mock.Anything).Return(role, nil)
	roleRepo.On("DeleteSoft", mock.Anything, mock.Anything).Return(errors.New("delete failed"))

	result := svc.DeleteRole(context.Background(), role.ID)

	assert.False(t, result.Success)
	assert.Equal(t, http.StatusInternalServerError, result.HttpStatus)
	roleRepo.AssertExpectations(t)
}

// ---------------------------------------------------------------------------
// AssignPermissions
// ---------------------------------------------------------------------------

func TestAssignPermissions_Success(t *testing.T) {
	roleRepo := new(MockRoleRepo)
	permRepo := new(MockPermissionRepo)
	svc := newTestService(roleRepo, permRepo)

	role := sampleRole()
	perm := samplePermission()
	permissionIDs := []uuid.UUID{perm.ID}

	roleRepo.On("FindOne", mock.Anything, mock.Anything, mock.Anything).Return(role, nil)
	permRepo.On("FindOne", mock.Anything, mock.Anything, mock.Anything).Return(perm, nil)
	roleRepo.On("AssignPermissions", mock.Anything, role.ID, mock.Anything).Return(nil)

	result := svc.AssignPermissions(context.Background(), role.ID, permissionIDs)

	assert.True(t, result.Success)
	assert.Equal(t, http.StatusOK, result.HttpStatus)
	assert.Equal(t, "Permissions assigned successfully", result.Message)
	roleRepo.AssertExpectations(t)
	permRepo.AssertExpectations(t)
}

func TestAssignPermissions_RoleNotFound(t *testing.T) {
	roleRepo := new(MockRoleRepo)
	svc := newTestService(roleRepo, new(MockPermissionRepo))

	roleRepo.On("FindOne", mock.Anything, mock.Anything, mock.Anything).Return(nil, gorm.ErrRecordNotFound)

	result := svc.AssignPermissions(context.Background(), uuid.New(), []uuid.UUID{uuid.New()})

	assert.False(t, result.Success)
	assert.Equal(t, http.StatusBadRequest, result.HttpStatus)
	assert.Equal(t, "Role not found", result.Message)
	roleRepo.AssertExpectations(t)
}

func TestAssignPermissions_PermissionNotFound(t *testing.T) {
	roleRepo := new(MockRoleRepo)
	permRepo := new(MockPermissionRepo)
	svc := newTestService(roleRepo, permRepo)

	role := sampleRole()
	roleRepo.On("FindOne", mock.Anything, mock.Anything, mock.Anything).Return(role, nil)
	permRepo.On("FindOne", mock.Anything, mock.Anything, mock.Anything).Return(nil, gorm.ErrRecordNotFound)

	result := svc.AssignPermissions(context.Background(), role.ID, []uuid.UUID{uuid.New()})

	assert.False(t, result.Success)
	assert.Equal(t, http.StatusBadRequest, result.HttpStatus)
	assert.Equal(t, "Permission not found", result.Message)
	roleRepo.AssertExpectations(t)
	permRepo.AssertExpectations(t)
}

// ---------------------------------------------------------------------------
// RemovePermission
// ---------------------------------------------------------------------------

func TestRemovePermission_Success(t *testing.T) {
	roleRepo := new(MockRoleRepo)
	permRepo := new(MockPermissionRepo)
	svc := newTestService(roleRepo, permRepo)

	role := sampleRole()
	perm := samplePermission()

	roleRepo.On("FindOne", mock.Anything, mock.Anything, mock.Anything).Return(role, nil)
	permRepo.On("FindOne", mock.Anything, mock.Anything, mock.Anything).Return(perm, nil)
	roleRepo.On("RemovePermission", mock.Anything, role.ID, perm).Return(nil)

	result := svc.RemovePermission(context.Background(), role.ID, perm.ID)

	assert.True(t, result.Success)
	assert.Equal(t, http.StatusOK, result.HttpStatus)
	assert.Equal(t, "Permission removed successfully", result.Message)
	roleRepo.AssertExpectations(t)
	permRepo.AssertExpectations(t)
}

func TestRemovePermission_RoleNotFound(t *testing.T) {
	roleRepo := new(MockRoleRepo)
	svc := newTestService(roleRepo, new(MockPermissionRepo))

	roleRepo.On("FindOne", mock.Anything, mock.Anything, mock.Anything).Return(nil, gorm.ErrRecordNotFound)

	result := svc.RemovePermission(context.Background(), uuid.New(), uuid.New())

	assert.False(t, result.Success)
	assert.Equal(t, http.StatusBadRequest, result.HttpStatus)
	assert.Equal(t, "Role not found", result.Message)
	roleRepo.AssertExpectations(t)
}
