package tagservice_test

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"manga-go/internal/pkg/common"
	"manga-go/internal/pkg/logger"
	"manga-go/internal/pkg/model"
	tagrequest "manga-go/internal/pkg/request/tag"
	tagservice "manga-go/internal/pkg/services/tag"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// MockTagRepo implements tagservice.TagRepository using testify/mock.
type MockTagRepo struct {
	mock.Mock
}

func (m *MockTagRepo) Create(ctx context.Context, tag *model.Tag) error {
	args := m.Called(ctx, tag)
	return args.Error(0)
}

func (m *MockTagRepo) FindOne(ctx context.Context, conditions []any, moreKeys map[string]common.MoreKeyOption) (*model.Tag, error) {
	args := m.Called(ctx, conditions, moreKeys)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Tag), args.Error(1)
}

func (m *MockTagRepo) DeleteSoft(ctx context.Context, conditions []any) error {
	args := m.Called(ctx, conditions)
	return args.Error(0)
}

func (m *MockTagRepo) FindPaginated(ctx context.Context, conditions []any, paging *common.Paging, moreKeys map[string]common.MoreKeyOption) ([]*model.Tag, int64, error) {
	args := m.Called(ctx, conditions, paging, moreKeys)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*model.Tag), args.Get(1).(int64), args.Error(2)
}

func (m *MockTagRepo) FindAll(ctx context.Context, conditions []any, moreKeys map[string]common.MoreKeyOption) ([]*model.Tag, error) {
	args := m.Called(ctx, conditions, moreKeys)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.Tag), args.Error(1)
}

func newTestService(repo *MockTagRepo) *tagservice.TagService {
	return tagservice.NewTagServiceWithRepo(logger.NewLogger(), repo)
}

func sampleTag() *model.Tag {
	now := time.Now()
	return &model.Tag{
		SqlModel: common.SqlModel{
			ID:        uuid.New(),
			CreatedAt: &now,
			UpdatedAt: &now,
		},
		Name: "Isekai",
		Slug: "isekai",
	}
}

// ---------------------------------------------------------------------------
// CreateTag
// ---------------------------------------------------------------------------

func TestCreateTag_Success(t *testing.T) {
	repo := new(MockTagRepo)
	svc := newTestService(repo)

	req := &tagrequest.CreateTagRequest{Name: "Isekai", Slug: "isekai"}
	repo.On("Create", mock.Anything, mock.AnythingOfType("*model.Tag")).Return(nil)

	result := svc.CreateTag(context.Background(), req)

	assert.True(t, result.Success)
	assert.Equal(t, http.StatusOK, result.HttpStatus)
	assert.Equal(t, "Tag created successfully", result.Message)
	repo.AssertExpectations(t)
}

func TestCreateTag_DBError(t *testing.T) {
	repo := new(MockTagRepo)
	svc := newTestService(repo)

	req := &tagrequest.CreateTagRequest{Name: "Isekai", Slug: "isekai"}
	repo.On("Create", mock.Anything, mock.AnythingOfType("*model.Tag")).Return(errors.New("db error"))

	result := svc.CreateTag(context.Background(), req)

	assert.False(t, result.Success)
	assert.Equal(t, http.StatusInternalServerError, result.HttpStatus)
	repo.AssertExpectations(t)
}

// ---------------------------------------------------------------------------
// ListTags
// ---------------------------------------------------------------------------

func TestListTags_Success(t *testing.T) {
	repo := new(MockTagRepo)
	svc := newTestService(repo)

	tags := []*model.Tag{sampleTag(), sampleTag()}
	paging := &common.Paging{Page: 1, Limit: 20}
	repo.On("FindPaginated", mock.Anything, mock.Anything, paging, mock.Anything).Return(tags, int64(2), nil)

	result := svc.ListTags(context.Background(), paging)

	assert.True(t, result.Success)
	assert.Equal(t, http.StatusOK, result.HttpStatus)
	repo.AssertExpectations(t)
}

func TestListTags_DBError(t *testing.T) {
	repo := new(MockTagRepo)
	svc := newTestService(repo)

	paging := &common.Paging{Page: 1, Limit: 20}
	repo.On("FindPaginated", mock.Anything, mock.Anything, paging, mock.Anything).Return(nil, int64(0), errors.New("db error"))

	result := svc.ListTags(context.Background(), paging)

	assert.False(t, result.Success)
	assert.Equal(t, http.StatusInternalServerError, result.HttpStatus)
	repo.AssertExpectations(t)
}

// ---------------------------------------------------------------------------
// ListAllTags
// ---------------------------------------------------------------------------

func TestListAllTags_Success(t *testing.T) {
	repo := new(MockTagRepo)
	svc := newTestService(repo)

	tags := []*model.Tag{sampleTag()}
	repo.On("FindAll", mock.Anything, mock.Anything, mock.Anything).Return(tags, nil)

	result := svc.ListAllTags(context.Background())

	assert.True(t, result.Success)
	assert.Equal(t, http.StatusOK, result.HttpStatus)
	repo.AssertExpectations(t)
}

func TestListAllTags_DBError(t *testing.T) {
	repo := new(MockTagRepo)
	svc := newTestService(repo)

	repo.On("FindAll", mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("db error"))

	result := svc.ListAllTags(context.Background())

	assert.False(t, result.Success)
	assert.Equal(t, http.StatusInternalServerError, result.HttpStatus)
	repo.AssertExpectations(t)
}

// ---------------------------------------------------------------------------
// DeleteTag
// ---------------------------------------------------------------------------

func TestDeleteTag_Success(t *testing.T) {
	repo := new(MockTagRepo)
	svc := newTestService(repo)

	tag := sampleTag()
	repo.On("FindOne", mock.Anything, mock.Anything, mock.Anything).Return(tag, nil)
	repo.On("DeleteSoft", mock.Anything, mock.Anything).Return(nil)

	result := svc.DeleteTag(context.Background(), "isekai")

	assert.True(t, result.Success)
	assert.Equal(t, http.StatusOK, result.HttpStatus)
	assert.Equal(t, "Tag deleted successfully", result.Message)
	repo.AssertExpectations(t)
}

func TestDeleteTag_NotFound(t *testing.T) {
	repo := new(MockTagRepo)
	svc := newTestService(repo)

	repo.On("FindOne", mock.Anything, mock.Anything, mock.Anything).Return(nil, gorm.ErrRecordNotFound)

	result := svc.DeleteTag(context.Background(), "nonexistent")

	assert.False(t, result.Success)
	assert.Equal(t, http.StatusBadRequest, result.HttpStatus)
	assert.Equal(t, "Tag not found", result.Message)
	repo.AssertExpectations(t)
}

func TestDeleteTag_DeleteDBError(t *testing.T) {
	repo := new(MockTagRepo)
	svc := newTestService(repo)

	tag := sampleTag()
	repo.On("FindOne", mock.Anything, mock.Anything, mock.Anything).Return(tag, nil)
	repo.On("DeleteSoft", mock.Anything, mock.Anything).Return(errors.New("delete failed"))

	result := svc.DeleteTag(context.Background(), "isekai")

	assert.False(t, result.Success)
	assert.Equal(t, http.StatusInternalServerError, result.HttpStatus)
	repo.AssertExpectations(t)
}
