package authorservice_test

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"manga-go/internal/pkg/common"
	"manga-go/internal/pkg/logger"
	"manga-go/internal/pkg/model"
	authorrequest "manga-go/internal/pkg/request/author"
	authorservice "manga-go/internal/pkg/services/author"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// MockAuthorRepo implements authorservice.AuthorRepository using testify/mock.
type MockAuthorRepo struct {
	mock.Mock
}

func (m *MockAuthorRepo) Create(ctx context.Context, author *model.Author) error {
	args := m.Called(ctx, author)
	return args.Error(0)
}

func (m *MockAuthorRepo) FindOne(ctx context.Context, conditions []any, moreKeys map[string]common.MoreKeyOption) (*model.Author, error) {
	args := m.Called(ctx, conditions, moreKeys)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Author), args.Error(1)
}

func (m *MockAuthorRepo) Update(ctx context.Context, conditions []any, data map[string]any) error {
	args := m.Called(ctx, conditions, data)
	return args.Error(0)
}

func (m *MockAuthorRepo) DeleteSoft(ctx context.Context, conditions []any) error {
	args := m.Called(ctx, conditions)
	return args.Error(0)
}

func (m *MockAuthorRepo) FindPaginated(ctx context.Context, conditions []any, paging *common.Paging, moreKeys map[string]common.MoreKeyOption) ([]*model.Author, int64, error) {
	args := m.Called(ctx, conditions, paging, moreKeys)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*model.Author), args.Get(1).(int64), args.Error(2)
}

func (m *MockAuthorRepo) FindAll(ctx context.Context, conditions []any, moreKeys map[string]common.MoreKeyOption) ([]*model.Author, error) {
	args := m.Called(ctx, conditions, moreKeys)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.Author), args.Error(1)
}

// newTestService creates an AuthorService with a mock repo for testing.
func newTestService(repo *MockAuthorRepo) *authorservice.AuthorService {
	return authorservice.NewAuthorServiceWithRepo(logger.NewLogger(), repo)
}

func sampleAuthor() *model.Author {
	now := time.Now()
	return &model.Author{
		SqlModel: common.SqlModel{
			ID:        uuid.New(),
			CreatedAt: &now,
			UpdatedAt: &now,
		},
		Name: "Test Author",
	}
}

// ---------------------------------------------------------------------------
// CreateAuthor
// ---------------------------------------------------------------------------

func TestCreateAuthor_Success(t *testing.T) {
	repo := new(MockAuthorRepo)
	svc := newTestService(repo)

	req := &authorrequest.CreateAuthorRequest{Name: "New Author"}
	repo.On("Create", mock.Anything, mock.AnythingOfType("*model.Author")).Return(nil)

	result := svc.CreateAuthor(context.Background(), req)

	assert.True(t, result.Success)
	assert.Equal(t, http.StatusOK, result.HttpStatus)
	assert.Equal(t, "Author created successfully", result.Message)
	repo.AssertExpectations(t)
}

func TestCreateAuthor_DBError(t *testing.T) {
	repo := new(MockAuthorRepo)
	svc := newTestService(repo)

	req := &authorrequest.CreateAuthorRequest{Name: "New Author"}
	dbErr := errors.New("db connection failed")
	repo.On("Create", mock.Anything, mock.AnythingOfType("*model.Author")).Return(dbErr)

	result := svc.CreateAuthor(context.Background(), req)

	assert.False(t, result.Success)
	assert.Equal(t, http.StatusInternalServerError, result.HttpStatus)
	repo.AssertExpectations(t)
}

// ---------------------------------------------------------------------------
// GetAuthor
// ---------------------------------------------------------------------------

func TestGetAuthor_Success(t *testing.T) {
	repo := new(MockAuthorRepo)
	svc := newTestService(repo)

	author := sampleAuthor()
	repo.On("FindOne", mock.Anything, mock.Anything, mock.Anything).Return(author, nil)

	result := svc.GetAuthor(context.Background(), author.ID)

	assert.True(t, result.Success)
	assert.Equal(t, http.StatusOK, result.HttpStatus)
	assert.Equal(t, "Author retrieved successfully", result.Message)
	repo.AssertExpectations(t)
}

func TestGetAuthor_NotFound(t *testing.T) {
	repo := new(MockAuthorRepo)
	svc := newTestService(repo)

	repo.On("FindOne", mock.Anything, mock.Anything, mock.Anything).Return(nil, gorm.ErrRecordNotFound)

	result := svc.GetAuthor(context.Background(), uuid.New())

	assert.False(t, result.Success)
	assert.Equal(t, http.StatusBadRequest, result.HttpStatus)
	assert.Equal(t, "Author not found", result.Message)
	repo.AssertExpectations(t)
}

func TestGetAuthor_DBError(t *testing.T) {
	repo := new(MockAuthorRepo)
	svc := newTestService(repo)

	repo.On("FindOne", mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("timeout"))

	result := svc.GetAuthor(context.Background(), uuid.New())

	assert.False(t, result.Success)
	assert.Equal(t, http.StatusInternalServerError, result.HttpStatus)
	repo.AssertExpectations(t)
}

// ---------------------------------------------------------------------------
// ListAuthors
// ---------------------------------------------------------------------------

func TestListAuthors_Success(t *testing.T) {
	repo := new(MockAuthorRepo)
	svc := newTestService(repo)

	authors := []*model.Author{sampleAuthor(), sampleAuthor()}
	// Each call to sampleAuthor() returns a struct with a freshly-generated UUID,
	// so the two elements in the slice are guaranteed to be distinct.
	paging := &common.Paging{Page: 1, Limit: 20}
	repo.On("FindPaginated", mock.Anything, mock.Anything, paging, mock.Anything).Return(authors, int64(2), nil)

	result := svc.ListAuthors(context.Background(), paging)

	assert.True(t, result.Success)
	assert.Equal(t, http.StatusOK, result.HttpStatus)
	repo.AssertExpectations(t)
}

func TestListAuthors_DBError(t *testing.T) {
	repo := new(MockAuthorRepo)
	svc := newTestService(repo)

	paging := &common.Paging{Page: 1, Limit: 20}
	repo.On("FindPaginated", mock.Anything, mock.Anything, paging, mock.Anything).Return(nil, int64(0), errors.New("db error"))

	result := svc.ListAuthors(context.Background(), paging)

	assert.False(t, result.Success)
	assert.Equal(t, http.StatusInternalServerError, result.HttpStatus)
	repo.AssertExpectations(t)
}

// ---------------------------------------------------------------------------
// ListAllAuthors
// ---------------------------------------------------------------------------

func TestListAllAuthors_Success(t *testing.T) {
	repo := new(MockAuthorRepo)
	svc := newTestService(repo)

	authors := []*model.Author{sampleAuthor()}
	repo.On("FindAll", mock.Anything, mock.Anything, mock.Anything).Return(authors, nil)

	result := svc.ListAllAuthors(context.Background())

	assert.True(t, result.Success)
	assert.Equal(t, http.StatusOK, result.HttpStatus)
	repo.AssertExpectations(t)
}

func TestListAllAuthors_DBError(t *testing.T) {
	repo := new(MockAuthorRepo)
	svc := newTestService(repo)

	repo.On("FindAll", mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("db error"))

	result := svc.ListAllAuthors(context.Background())

	assert.False(t, result.Success)
	assert.Equal(t, http.StatusInternalServerError, result.HttpStatus)
	repo.AssertExpectations(t)
}

// ---------------------------------------------------------------------------
// UpdateAuthor
// ---------------------------------------------------------------------------

func TestUpdateAuthor_Success(t *testing.T) {
	repo := new(MockAuthorRepo)
	svc := newTestService(repo)

	author := sampleAuthor()
	req := &authorrequest.UpdateAuthorRequest{Name: "Updated Name"}

	repo.On("FindOne", mock.Anything, mock.Anything, mock.Anything).Return(author, nil)
	repo.On("Update", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	result := svc.UpdateAuthor(context.Background(), author.ID, req)

	assert.True(t, result.Success)
	assert.Equal(t, http.StatusOK, result.HttpStatus)
	assert.Equal(t, "Author updated successfully", result.Message)
	repo.AssertExpectations(t)
}

func TestUpdateAuthor_NotFound(t *testing.T) {
	repo := new(MockAuthorRepo)
	svc := newTestService(repo)

	req := &authorrequest.UpdateAuthorRequest{Name: "Updated Name"}
	repo.On("FindOne", mock.Anything, mock.Anything, mock.Anything).Return(nil, gorm.ErrRecordNotFound)

	result := svc.UpdateAuthor(context.Background(), uuid.New(), req)

	assert.False(t, result.Success)
	assert.Equal(t, http.StatusBadRequest, result.HttpStatus)
	assert.Equal(t, "Author not found", result.Message)
	repo.AssertExpectations(t)
}

func TestUpdateAuthor_UpdateDBError(t *testing.T) {
	repo := new(MockAuthorRepo)
	svc := newTestService(repo)

	author := sampleAuthor()
	req := &authorrequest.UpdateAuthorRequest{Name: "Updated Name"}

	repo.On("FindOne", mock.Anything, mock.Anything, mock.Anything).Return(author, nil)
	repo.On("Update", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("update failed"))

	result := svc.UpdateAuthor(context.Background(), author.ID, req)

	assert.False(t, result.Success)
	assert.Equal(t, http.StatusInternalServerError, result.HttpStatus)
	repo.AssertExpectations(t)
}

// ---------------------------------------------------------------------------
// DeleteAuthor
// ---------------------------------------------------------------------------

func TestDeleteAuthor_Success(t *testing.T) {
	repo := new(MockAuthorRepo)
	svc := newTestService(repo)

	author := sampleAuthor()
	repo.On("FindOne", mock.Anything, mock.Anything, mock.Anything).Return(author, nil)
	repo.On("DeleteSoft", mock.Anything, mock.Anything).Return(nil)

	result := svc.DeleteAuthor(context.Background(), author.ID)

	assert.True(t, result.Success)
	assert.Equal(t, http.StatusOK, result.HttpStatus)
	assert.Equal(t, "Author deleted successfully", result.Message)
	repo.AssertExpectations(t)
}

func TestDeleteAuthor_NotFound(t *testing.T) {
	repo := new(MockAuthorRepo)
	svc := newTestService(repo)

	repo.On("FindOne", mock.Anything, mock.Anything, mock.Anything).Return(nil, gorm.ErrRecordNotFound)

	result := svc.DeleteAuthor(context.Background(), uuid.New())

	assert.False(t, result.Success)
	assert.Equal(t, http.StatusBadRequest, result.HttpStatus)
	assert.Equal(t, "Author not found", result.Message)
	repo.AssertExpectations(t)
}

func TestDeleteAuthor_DeleteDBError(t *testing.T) {
	repo := new(MockAuthorRepo)
	svc := newTestService(repo)

	author := sampleAuthor()
	repo.On("FindOne", mock.Anything, mock.Anything, mock.Anything).Return(author, nil)
	repo.On("DeleteSoft", mock.Anything, mock.Anything).Return(errors.New("delete failed"))

	result := svc.DeleteAuthor(context.Background(), author.ID)

	assert.False(t, result.Success)
	assert.Equal(t, http.StatusInternalServerError, result.HttpStatus)
	repo.AssertExpectations(t)
}
