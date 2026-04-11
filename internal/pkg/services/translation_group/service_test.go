package translationgroupservice_test

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"manga-go/internal/pkg/common"
	"manga-go/internal/pkg/logger"
	"manga-go/internal/pkg/model"
	translationgrouprequest "manga-go/internal/pkg/request/translation_group"
	translationgroupservice "manga-go/internal/pkg/services/translation_group"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// MockTranslationGroupRepo implements translationgroupservice.TranslationGroupRepository.
type MockTranslationGroupRepo struct {
	mock.Mock
}

func (m *MockTranslationGroupRepo) Create(ctx context.Context, group *model.TranslationGroup) error {
	args := m.Called(ctx, group)
	return args.Error(0)
}

func (m *MockTranslationGroupRepo) FindOne(ctx context.Context, conditions []any, moreKeys map[string]common.MoreKeyOption) (*model.TranslationGroup, error) {
	args := m.Called(ctx, conditions, moreKeys)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.TranslationGroup), args.Error(1)
}

func (m *MockTranslationGroupRepo) Update(ctx context.Context, conditions []any, data map[string]any) error {
	args := m.Called(ctx, conditions, data)
	return args.Error(0)
}

func (m *MockTranslationGroupRepo) DeleteSoft(ctx context.Context, conditions []any) error {
	args := m.Called(ctx, conditions)
	return args.Error(0)
}

func (m *MockTranslationGroupRepo) FindPaginated(ctx context.Context, conditions []any, paging *common.Paging, moreKeys map[string]common.MoreKeyOption) ([]*model.TranslationGroup, int64, error) {
	args := m.Called(ctx, conditions, paging, moreKeys)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*model.TranslationGroup), args.Get(1).(int64), args.Error(2)
}

// MockUserRepo implements translationgroupservice.UserRepository.
type MockUserRepo struct {
	mock.Mock
}

func (m *MockUserRepo) FindOne(ctx context.Context, conditions []any, moreKeys map[string]common.MoreKeyOption) (*model.User, error) {
	args := m.Called(ctx, conditions, moreKeys)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepo) Update(ctx context.Context, conditions []any, data map[string]any) error {
	args := m.Called(ctx, conditions, data)
	return args.Error(0)
}

func newTestService(tgRepo *MockTranslationGroupRepo, userRepo *MockUserRepo) *translationgroupservice.TranslationGroupService {
	return translationgroupservice.NewTranslationGroupServiceWithRepos(logger.NewLogger(), tgRepo, userRepo)
}

func sampleGroup(ownerID uuid.UUID) *model.TranslationGroup {
	now := time.Now()
	return &model.TranslationGroup{
		SqlModel: common.SqlModel{
			ID:        uuid.New(),
			CreatedAt: &now,
			UpdatedAt: &now,
		},
		Name:    "Cool Scanners",
		Slug:    "cool-scanners",
		OwnerID: ownerID,
	}
}

func sampleUser() *model.User {
	now := time.Now()
	return &model.User{
		SqlModel: common.SqlModel{
			ID:        uuid.New(),
			CreatedAt: &now,
			UpdatedAt: &now,
		},
		Name:  "Test User",
		Email: "test@example.com",
	}
}

// ---------------------------------------------------------------------------
// CreateTranslationGroup
// ---------------------------------------------------------------------------

func TestCreateTranslationGroup_Success(t *testing.T) {
	tgRepo := new(MockTranslationGroupRepo)
	userRepo := new(MockUserRepo)
	svc := newTestService(tgRepo, userRepo)

	ownerID := uuid.New()
	req := &translationgrouprequest.CreateTranslationGroupRequest{Name: "Cool Scanners", Slug: "cool-scanners"}

	tgRepo.On("Create", mock.Anything, mock.AnythingOfType("*model.TranslationGroup")).Return(nil)
	userRepo.On("Update", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	result := svc.CreateTranslationGroup(context.Background(), ownerID, req)

	assert.True(t, result.Success)
	assert.Equal(t, http.StatusOK, result.HttpStatus)
	assert.Equal(t, "Translation group created successfully", result.Message)
	tgRepo.AssertExpectations(t)
	userRepo.AssertExpectations(t)
}

func TestCreateTranslationGroup_DBError(t *testing.T) {
	tgRepo := new(MockTranslationGroupRepo)
	svc := newTestService(tgRepo, new(MockUserRepo))

	ownerID := uuid.New()
	req := &translationgrouprequest.CreateTranslationGroupRequest{Name: "Cool Scanners", Slug: "cool-scanners"}

	tgRepo.On("Create", mock.Anything, mock.AnythingOfType("*model.TranslationGroup")).Return(errors.New("db error"))

	result := svc.CreateTranslationGroup(context.Background(), ownerID, req)

	assert.False(t, result.Success)
	assert.Equal(t, http.StatusInternalServerError, result.HttpStatus)
	tgRepo.AssertExpectations(t)
}

// ---------------------------------------------------------------------------
// GetTranslationGroup
// ---------------------------------------------------------------------------

func TestGetTranslationGroup_Success(t *testing.T) {
	tgRepo := new(MockTranslationGroupRepo)
	svc := newTestService(tgRepo, new(MockUserRepo))

	group := sampleGroup(uuid.New())
	tgRepo.On("FindOne", mock.Anything, mock.Anything, mock.Anything).Return(group, nil)

	result := svc.GetTranslationGroup(context.Background(), group.Slug)

	assert.True(t, result.Success)
	assert.Equal(t, http.StatusOK, result.HttpStatus)
	assert.Equal(t, "Translation group retrieved successfully", result.Message)
	tgRepo.AssertExpectations(t)
}

func TestGetTranslationGroup_NotFound(t *testing.T) {
	tgRepo := new(MockTranslationGroupRepo)
	svc := newTestService(tgRepo, new(MockUserRepo))

	tgRepo.On("FindOne", mock.Anything, mock.Anything, mock.Anything).Return(nil, gorm.ErrRecordNotFound)

	result := svc.GetTranslationGroup(context.Background(), "missing-slug")

	assert.False(t, result.Success)
	assert.Equal(t, http.StatusBadRequest, result.HttpStatus)
	assert.Equal(t, "TranslationGroup not found", result.Message)
	tgRepo.AssertExpectations(t)
}

func TestGetTranslationGroup_DBError(t *testing.T) {
	tgRepo := new(MockTranslationGroupRepo)
	svc := newTestService(tgRepo, new(MockUserRepo))

	tgRepo.On("FindOne", mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("timeout"))

	result := svc.GetTranslationGroup(context.Background(), "slug")

	assert.False(t, result.Success)
	assert.Equal(t, http.StatusInternalServerError, result.HttpStatus)
	tgRepo.AssertExpectations(t)
}

// ---------------------------------------------------------------------------
// ListTranslationGroups
// ---------------------------------------------------------------------------

func TestListTranslationGroups_Success(t *testing.T) {
	tgRepo := new(MockTranslationGroupRepo)
	svc := newTestService(tgRepo, new(MockUserRepo))

	groups := []*model.TranslationGroup{sampleGroup(uuid.New()), sampleGroup(uuid.New())}
	paging := &common.Paging{Page: 1, Limit: 20}
	tgRepo.On("FindPaginated", mock.Anything, mock.Anything, paging, mock.Anything).Return(groups, int64(2), nil)

	result := svc.ListTranslationGroups(context.Background(), paging)

	assert.True(t, result.Success)
	assert.Equal(t, http.StatusOK, result.HttpStatus)
	tgRepo.AssertExpectations(t)
}

func TestListTranslationGroups_DBError(t *testing.T) {
	tgRepo := new(MockTranslationGroupRepo)
	svc := newTestService(tgRepo, new(MockUserRepo))

	paging := &common.Paging{Page: 1, Limit: 20}
	tgRepo.On("FindPaginated", mock.Anything, mock.Anything, paging, mock.Anything).Return(nil, int64(0), errors.New("db error"))

	result := svc.ListTranslationGroups(context.Background(), paging)

	assert.False(t, result.Success)
	assert.Equal(t, http.StatusInternalServerError, result.HttpStatus)
	tgRepo.AssertExpectations(t)
}

// ---------------------------------------------------------------------------
// UpdateTranslationGroup
// ---------------------------------------------------------------------------

func TestUpdateTranslationGroup_Success(t *testing.T) {
	tgRepo := new(MockTranslationGroupRepo)
	svc := newTestService(tgRepo, new(MockUserRepo))

	group := sampleGroup(uuid.New())
	req := &translationgrouprequest.UpdateTranslationGroupRequest{Name: "Elite Scanners", Slug: "elite-scanners"}

	tgRepo.On("FindOne", mock.Anything, mock.Anything, mock.Anything).Return(group, nil)
	tgRepo.On("Update", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	result := svc.UpdateTranslationGroup(context.Background(), group.Slug, req)

	assert.True(t, result.Success)
	assert.Equal(t, http.StatusOK, result.HttpStatus)
	assert.Equal(t, "Translation group updated successfully", result.Message)
	tgRepo.AssertExpectations(t)
}

func TestUpdateTranslationGroup_NotFound(t *testing.T) {
	tgRepo := new(MockTranslationGroupRepo)
	svc := newTestService(tgRepo, new(MockUserRepo))

	req := &translationgrouprequest.UpdateTranslationGroupRequest{Name: "Elite Scanners", Slug: "elite-scanners"}
	tgRepo.On("FindOne", mock.Anything, mock.Anything, mock.Anything).Return(nil, gorm.ErrRecordNotFound)

	result := svc.UpdateTranslationGroup(context.Background(), "missing-slug", req)

	assert.False(t, result.Success)
	assert.Equal(t, http.StatusBadRequest, result.HttpStatus)
	assert.Equal(t, "TranslationGroup not found", result.Message)
	tgRepo.AssertExpectations(t)
}

func TestUpdateTranslationGroup_UpdateDBError(t *testing.T) {
	tgRepo := new(MockTranslationGroupRepo)
	svc := newTestService(tgRepo, new(MockUserRepo))

	group := sampleGroup(uuid.New())
	req := &translationgrouprequest.UpdateTranslationGroupRequest{Name: "Elite Scanners", Slug: "elite-scanners"}

	tgRepo.On("FindOne", mock.Anything, mock.Anything, mock.Anything).Return(group, nil)
	tgRepo.On("Update", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("update failed"))

	result := svc.UpdateTranslationGroup(context.Background(), group.Slug, req)

	assert.False(t, result.Success)
	assert.Equal(t, http.StatusInternalServerError, result.HttpStatus)
	tgRepo.AssertExpectations(t)
}

// ---------------------------------------------------------------------------
// DeleteTranslationGroup
// ---------------------------------------------------------------------------

func TestDeleteTranslationGroup_Success(t *testing.T) {
	tgRepo := new(MockTranslationGroupRepo)
	svc := newTestService(tgRepo, new(MockUserRepo))

	group := sampleGroup(uuid.New())
	tgRepo.On("FindOne", mock.Anything, mock.Anything, mock.Anything).Return(group, nil)
	tgRepo.On("DeleteSoft", mock.Anything, mock.Anything).Return(nil)

	result := svc.DeleteTranslationGroup(context.Background(), group.Slug)

	assert.True(t, result.Success)
	assert.Equal(t, http.StatusOK, result.HttpStatus)
	assert.Equal(t, "Translation group deleted successfully", result.Message)
	tgRepo.AssertExpectations(t)
}

func TestDeleteTranslationGroup_NotFound(t *testing.T) {
	tgRepo := new(MockTranslationGroupRepo)
	svc := newTestService(tgRepo, new(MockUserRepo))

	tgRepo.On("FindOne", mock.Anything, mock.Anything, mock.Anything).Return(nil, gorm.ErrRecordNotFound)

	result := svc.DeleteTranslationGroup(context.Background(), "missing-slug")

	assert.False(t, result.Success)
	assert.Equal(t, http.StatusBadRequest, result.HttpStatus)
	assert.Equal(t, "TranslationGroup not found", result.Message)
	tgRepo.AssertExpectations(t)
}

func TestDeleteTranslationGroup_DeleteDBError(t *testing.T) {
	tgRepo := new(MockTranslationGroupRepo)
	svc := newTestService(tgRepo, new(MockUserRepo))

	group := sampleGroup(uuid.New())
	tgRepo.On("FindOne", mock.Anything, mock.Anything, mock.Anything).Return(group, nil)
	tgRepo.On("DeleteSoft", mock.Anything, mock.Anything).Return(errors.New("delete failed"))

	result := svc.DeleteTranslationGroup(context.Background(), group.Slug)

	assert.False(t, result.Success)
	assert.Equal(t, http.StatusInternalServerError, result.HttpStatus)
	tgRepo.AssertExpectations(t)
}

// ---------------------------------------------------------------------------
// TransferOwnership
// ---------------------------------------------------------------------------

func TestTransferOwnership_Success(t *testing.T) {
	tgRepo := new(MockTranslationGroupRepo)
	userRepo := new(MockUserRepo)
	svc := newTestService(tgRepo, userRepo)

	group := sampleGroup(uuid.New())
	newOwner := sampleUser()
	req := &translationgrouprequest.TransferOwnershipRequest{NewOwnerID: newOwner.ID}

	tgRepo.On("FindOne", mock.Anything, mock.Anything, mock.Anything).Return(group, nil)
	userRepo.On("FindOne", mock.Anything, mock.Anything, mock.Anything).Return(newOwner, nil)
	tgRepo.On("Update", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	result := svc.TransferOwnership(context.Background(), group.Slug, req)

	assert.True(t, result.Success)
	assert.Equal(t, http.StatusOK, result.HttpStatus)
	assert.Equal(t, "Ownership transferred successfully", result.Message)
	tgRepo.AssertExpectations(t)
	userRepo.AssertExpectations(t)
}

func TestTransferOwnership_NewOwnerNotMember(t *testing.T) {
	tgRepo := new(MockTranslationGroupRepo)
	userRepo := new(MockUserRepo)
	svc := newTestService(tgRepo, userRepo)

	group := sampleGroup(uuid.New())
	req := &translationgrouprequest.TransferOwnershipRequest{NewOwnerID: uuid.New()}

	tgRepo.On("FindOne", mock.Anything, mock.Anything, mock.Anything).Return(group, nil)
	userRepo.On("FindOne", mock.Anything, mock.Anything, mock.Anything).Return(nil, gorm.ErrRecordNotFound)

	result := svc.TransferOwnership(context.Background(), group.Slug, req)

	assert.False(t, result.Success)
	assert.Equal(t, http.StatusBadRequest, result.HttpStatus)
	assert.Equal(t, "New owner must be a member of the translation group", result.Message)
	tgRepo.AssertExpectations(t)
	userRepo.AssertExpectations(t)
}

func TestTransferOwnership_GroupNotFound(t *testing.T) {
	tgRepo := new(MockTranslationGroupRepo)
	svc := newTestService(tgRepo, new(MockUserRepo))

	tgRepo.On("FindOne", mock.Anything, mock.Anything, mock.Anything).Return(nil, gorm.ErrRecordNotFound)

	result := svc.TransferOwnership(context.Background(), "missing-slug", &translationgrouprequest.TransferOwnershipRequest{NewOwnerID: uuid.New()})

	assert.False(t, result.Success)
	assert.Equal(t, http.StatusBadRequest, result.HttpStatus)
	assert.Equal(t, "TranslationGroup not found", result.Message)
	tgRepo.AssertExpectations(t)
}
