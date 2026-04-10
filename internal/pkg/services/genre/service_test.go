package genreservice_test

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"manga-go/internal/pkg/common"
	"manga-go/internal/pkg/logger"
	"manga-go/internal/pkg/model"
	genrerequest "manga-go/internal/pkg/request/genre"
	genreservice "manga-go/internal/pkg/services/genre"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// MockGenreRepo implements genreservice.GenreRepository using testify/mock.
type MockGenreRepo struct {
	mock.Mock
}

func (m *MockGenreRepo) Create(ctx context.Context, genre *model.Genre) error {
	args := m.Called(ctx, genre)
	return args.Error(0)
}

func (m *MockGenreRepo) FindOne(ctx context.Context, conditions []any, moreKeys map[string]common.MoreKeyOption) (*model.Genre, error) {
	args := m.Called(ctx, conditions, moreKeys)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Genre), args.Error(1)
}

func (m *MockGenreRepo) Update(ctx context.Context, conditions []any, data map[string]any) error {
	args := m.Called(ctx, conditions, data)
	return args.Error(0)
}

func (m *MockGenreRepo) DeleteSoft(ctx context.Context, conditions []any) error {
	args := m.Called(ctx, conditions)
	return args.Error(0)
}

func (m *MockGenreRepo) FindPaginated(ctx context.Context, conditions []any, paging *common.Paging, moreKeys map[string]common.MoreKeyOption) ([]*model.Genre, int64, error) {
	args := m.Called(ctx, conditions, paging, moreKeys)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*model.Genre), args.Get(1).(int64), args.Error(2)
}

func (m *MockGenreRepo) FindAll(ctx context.Context, conditions []any, moreKeys map[string]common.MoreKeyOption) ([]*model.Genre, error) {
	args := m.Called(ctx, conditions, moreKeys)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.Genre), args.Error(1)
}

func newTestService(repo *MockGenreRepo) *genreservice.GenreService {
	return genreservice.NewGenreServiceWithRepo(logger.NewLogger(), repo)
}

func sampleGenre() *model.Genre {
	now := time.Now()
	return &model.Genre{
		SqlModel: common.SqlModel{
			ID:        uuid.New(),
			CreatedAt: &now,
			UpdatedAt: &now,
		},
		Name:        "Action",
		Slug:        "action",
		Description: "Action manga",
		Thumbnail:   "action.jpg",
	}
}

// ---------------------------------------------------------------------------
// CreateGenre
// ---------------------------------------------------------------------------

func TestCreateGenre_Success(t *testing.T) {
	repo := new(MockGenreRepo)
	svc := newTestService(repo)

	req := &genrerequest.CreateGenreRequest{
		Name:        "Action",
		Slug:        "action",
		Description: "Action manga",
		Thumbnail:   "action.jpg",
	}
	repo.On("Create", mock.Anything, mock.AnythingOfType("*model.Genre")).Return(nil)

	result := svc.CreateGenre(context.Background(), req)

	assert.True(t, result.Success)
	assert.Equal(t, http.StatusOK, result.HttpStatus)
	assert.Equal(t, "Genre created successfully", result.Message)
	repo.AssertExpectations(t)
}

func TestCreateGenre_DBError(t *testing.T) {
	repo := new(MockGenreRepo)
	svc := newTestService(repo)

	req := &genrerequest.CreateGenreRequest{Name: "Action", Slug: "action"}
	repo.On("Create", mock.Anything, mock.AnythingOfType("*model.Genre")).Return(errors.New("db error"))

	result := svc.CreateGenre(context.Background(), req)

	assert.False(t, result.Success)
	assert.Equal(t, http.StatusInternalServerError, result.HttpStatus)
	repo.AssertExpectations(t)
}

// ---------------------------------------------------------------------------
// GetGenre
// ---------------------------------------------------------------------------

func TestGetGenre_Success(t *testing.T) {
	repo := new(MockGenreRepo)
	svc := newTestService(repo)

	genre := sampleGenre()
	repo.On("FindOne", mock.Anything, mock.Anything, mock.Anything).Return(genre, nil)

	result := svc.GetGenre(context.Background(), "action")

	assert.True(t, result.Success)
	assert.Equal(t, http.StatusOK, result.HttpStatus)
	assert.Equal(t, "Genre retrieved successfully", result.Message)
	repo.AssertExpectations(t)
}

func TestGetGenre_NotFound(t *testing.T) {
	repo := new(MockGenreRepo)
	svc := newTestService(repo)

	repo.On("FindOne", mock.Anything, mock.Anything, mock.Anything).Return(nil, gorm.ErrRecordNotFound)

	result := svc.GetGenre(context.Background(), "nonexistent")

	assert.False(t, result.Success)
	assert.Equal(t, http.StatusBadRequest, result.HttpStatus)
	assert.Equal(t, "Genre not found", result.Message)
	repo.AssertExpectations(t)
}

func TestGetGenre_DBError(t *testing.T) {
	repo := new(MockGenreRepo)
	svc := newTestService(repo)

	repo.On("FindOne", mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("timeout"))

	result := svc.GetGenre(context.Background(), "action")

	assert.False(t, result.Success)
	assert.Equal(t, http.StatusInternalServerError, result.HttpStatus)
	repo.AssertExpectations(t)
}

// ---------------------------------------------------------------------------
// ListGenres
// ---------------------------------------------------------------------------

func TestListGenres_Success(t *testing.T) {
	repo := new(MockGenreRepo)
	svc := newTestService(repo)

	genres := []*model.Genre{sampleGenre(), sampleGenre()}
	// Each call to sampleGenre() returns a struct with a freshly-generated UUID,
	// so the two elements in the slice are guaranteed to be distinct.
	paging := &common.Paging{Page: 1, Limit: 20}
	repo.On("FindPaginated", mock.Anything, mock.Anything, paging, mock.Anything).Return(genres, int64(2), nil)

	result := svc.ListGenres(context.Background(), paging)

	assert.True(t, result.Success)
	assert.Equal(t, http.StatusOK, result.HttpStatus)
	repo.AssertExpectations(t)
}

func TestListGenres_DBError(t *testing.T) {
	repo := new(MockGenreRepo)
	svc := newTestService(repo)

	paging := &common.Paging{Page: 1, Limit: 20}
	repo.On("FindPaginated", mock.Anything, mock.Anything, paging, mock.Anything).Return(nil, int64(0), errors.New("db error"))

	result := svc.ListGenres(context.Background(), paging)

	assert.False(t, result.Success)
	assert.Equal(t, http.StatusInternalServerError, result.HttpStatus)
	repo.AssertExpectations(t)
}

// ---------------------------------------------------------------------------
// ListAllGenres
// ---------------------------------------------------------------------------

func TestListAllGenres_Success(t *testing.T) {
	repo := new(MockGenreRepo)
	svc := newTestService(repo)

	genres := []*model.Genre{sampleGenre()}
	repo.On("FindAll", mock.Anything, mock.Anything, mock.Anything).Return(genres, nil)

	result := svc.ListAllGenres(context.Background())

	assert.True(t, result.Success)
	assert.Equal(t, http.StatusOK, result.HttpStatus)
	repo.AssertExpectations(t)
}

func TestListAllGenres_DBError(t *testing.T) {
	repo := new(MockGenreRepo)
	svc := newTestService(repo)

	repo.On("FindAll", mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("db error"))

	result := svc.ListAllGenres(context.Background())

	assert.False(t, result.Success)
	assert.Equal(t, http.StatusInternalServerError, result.HttpStatus)
	repo.AssertExpectations(t)
}

// ---------------------------------------------------------------------------
// UpdateGenre
// ---------------------------------------------------------------------------

func TestUpdateGenre_Success(t *testing.T) {
	repo := new(MockGenreRepo)
	svc := newTestService(repo)

	genre := sampleGenre()
	req := &genrerequest.UpdateGenreRequest{
		Name:        "Updated Action",
		Slug:        "updated-action",
		Description: "Updated description",
		Thumbnail:   "new.jpg",
	}

	repo.On("FindOne", mock.Anything, mock.Anything, mock.Anything).Return(genre, nil)
	repo.On("Update", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	result := svc.UpdateGenre(context.Background(), "action", req)

	assert.True(t, result.Success)
	assert.Equal(t, http.StatusOK, result.HttpStatus)
	assert.Equal(t, "Genre updated successfully", result.Message)
	repo.AssertExpectations(t)
}

func TestUpdateGenre_NotFound(t *testing.T) {
	repo := new(MockGenreRepo)
	svc := newTestService(repo)

	req := &genrerequest.UpdateGenreRequest{Name: "Updated", Slug: "updated"}
	repo.On("FindOne", mock.Anything, mock.Anything, mock.Anything).Return(nil, gorm.ErrRecordNotFound)

	result := svc.UpdateGenre(context.Background(), "nonexistent", req)

	assert.False(t, result.Success)
	assert.Equal(t, http.StatusBadRequest, result.HttpStatus)
	assert.Equal(t, "Genre not found", result.Message)
	repo.AssertExpectations(t)
}

func TestUpdateGenre_UpdateDBError(t *testing.T) {
	repo := new(MockGenreRepo)
	svc := newTestService(repo)

	genre := sampleGenre()
	req := &genrerequest.UpdateGenreRequest{Name: "Updated", Slug: "updated"}

	repo.On("FindOne", mock.Anything, mock.Anything, mock.Anything).Return(genre, nil)
	repo.On("Update", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("update failed"))

	result := svc.UpdateGenre(context.Background(), "action", req)

	assert.False(t, result.Success)
	assert.Equal(t, http.StatusInternalServerError, result.HttpStatus)
	repo.AssertExpectations(t)
}

// ---------------------------------------------------------------------------
// DeleteGenre
// ---------------------------------------------------------------------------

func TestDeleteGenre_Success(t *testing.T) {
	repo := new(MockGenreRepo)
	svc := newTestService(repo)

	genre := sampleGenre()
	repo.On("FindOne", mock.Anything, mock.Anything, mock.Anything).Return(genre, nil)
	repo.On("DeleteSoft", mock.Anything, mock.Anything).Return(nil)

	result := svc.DeleteGenre(context.Background(), "action")

	assert.True(t, result.Success)
	assert.Equal(t, http.StatusOK, result.HttpStatus)
	assert.Equal(t, "Genre deleted successfully", result.Message)
	repo.AssertExpectations(t)
}

func TestDeleteGenre_NotFound(t *testing.T) {
	repo := new(MockGenreRepo)
	svc := newTestService(repo)

	repo.On("FindOne", mock.Anything, mock.Anything, mock.Anything).Return(nil, gorm.ErrRecordNotFound)

	result := svc.DeleteGenre(context.Background(), "nonexistent")

	assert.False(t, result.Success)
	assert.Equal(t, http.StatusBadRequest, result.HttpStatus)
	assert.Equal(t, "Genre not found", result.Message)
	repo.AssertExpectations(t)
}

func TestDeleteGenre_DeleteDBError(t *testing.T) {
	repo := new(MockGenreRepo)
	svc := newTestService(repo)

	genre := sampleGenre()
	repo.On("FindOne", mock.Anything, mock.Anything, mock.Anything).Return(genre, nil)
	repo.On("DeleteSoft", mock.Anything, mock.Anything).Return(errors.New("delete failed"))

	result := svc.DeleteGenre(context.Background(), "action")

	assert.False(t, result.Success)
	assert.Equal(t, http.StatusInternalServerError, result.HttpStatus)
	repo.AssertExpectations(t)
}
