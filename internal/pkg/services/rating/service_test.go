package ratingservice_test

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"manga-go/internal/pkg/common"
	"manga-go/internal/pkg/logger"
	"manga-go/internal/pkg/model"
	ratingrequest "manga-go/internal/pkg/request/rating"
	ratingservice "manga-go/internal/pkg/services/rating"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// MockRatingRepo implements ratingservice.RatingRepository using testify/mock.
type MockRatingRepo struct {
	mock.Mock
}

func (m *MockRatingRepo) Create(ctx context.Context, rating *model.Rating) error {
	args := m.Called(ctx, rating)
	return args.Error(0)
}

func (m *MockRatingRepo) FindOne(ctx context.Context, conditions []any, moreKeys map[string]common.MoreKeyOption) (*model.Rating, error) {
	args := m.Called(ctx, conditions, moreKeys)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Rating), args.Error(1)
}

func (m *MockRatingRepo) Update(ctx context.Context, conditions []any, data map[string]any) error {
	args := m.Called(ctx, conditions, data)
	return args.Error(0)
}

func (m *MockRatingRepo) DeleteSoft(ctx context.Context, conditions []any) error {
	args := m.Called(ctx, conditions)
	return args.Error(0)
}

func (m *MockRatingRepo) FindPaginated(ctx context.Context, conditions []any, paging *common.Paging, moreKeys map[string]common.MoreKeyOption) ([]*model.Rating, int64, error) {
	args := m.Called(ctx, conditions, paging, moreKeys)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*model.Rating), args.Get(1).(int64), args.Error(2)
}

func (m *MockRatingRepo) FindByUserAndComic(ctx context.Context, userID uuid.UUID, comicID uuid.UUID) (*model.Rating, error) {
	args := m.Called(ctx, userID, comicID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Rating), args.Error(1)
}

func (m *MockRatingRepo) GetAverageRatingByComicID(ctx context.Context, comicID uuid.UUID) (float64, int64, error) {
	args := m.Called(ctx, comicID)
	return args.Get(0).(float64), args.Get(1).(int64), args.Error(2)
}

func newTestService(repo *MockRatingRepo) *ratingservice.RatingService {
	return ratingservice.NewRatingServiceWithRepo(logger.NewLogger(), repo)
}

func sampleRating(userID, comicID uuid.UUID) *model.Rating {
	now := time.Now()
	return &model.Rating{
		SqlModel: common.SqlModel{
			ID:        uuid.New(),
			CreatedAt: &now,
			UpdatedAt: &now,
		},
		UserId:  userID,
		ComicId: comicID,
		Score:   4,
	}
}

// ---------------------------------------------------------------------------
// CreateRating — no existing rating (creates new)
// ---------------------------------------------------------------------------

func TestCreateRating_NewRating_Success(t *testing.T) {
	repo := new(MockRatingRepo)
	svc := newTestService(repo)

	userID := uuid.New()
	comicID := uuid.New()
	req := &ratingrequest.CreateRatingRequest{Score: 4}

	repo.On("FindByUserAndComic", mock.Anything, userID, comicID).Return(nil, gorm.ErrRecordNotFound)
	repo.On("Create", mock.Anything, mock.AnythingOfType("*model.Rating")).Return(nil)

	result := svc.CreateRating(context.Background(), userID, comicID, req)

	assert.True(t, result.Success)
	assert.Equal(t, http.StatusOK, result.HttpStatus)
	assert.Equal(t, "Rating created successfully", result.Message)
	repo.AssertExpectations(t)
}

// CreateRating — existing rating (updates it)
func TestCreateRating_UpdateExisting_Success(t *testing.T) {
	repo := new(MockRatingRepo)
	svc := newTestService(repo)

	userID := uuid.New()
	comicID := uuid.New()
	existing := sampleRating(userID, comicID)
	req := &ratingrequest.CreateRatingRequest{Score: 5}

	repo.On("FindByUserAndComic", mock.Anything, userID, comicID).Return(existing, nil)
	repo.On("Update", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	result := svc.CreateRating(context.Background(), userID, comicID, req)

	assert.True(t, result.Success)
	assert.Equal(t, http.StatusOK, result.HttpStatus)
	assert.Equal(t, "Rating updated successfully", result.Message)
	repo.AssertExpectations(t)
}

func TestCreateRating_DBError(t *testing.T) {
	repo := new(MockRatingRepo)
	svc := newTestService(repo)

	userID := uuid.New()
	comicID := uuid.New()
	req := &ratingrequest.CreateRatingRequest{Score: 3}

	repo.On("FindByUserAndComic", mock.Anything, userID, comicID).Return(nil, errors.New("db error"))

	result := svc.CreateRating(context.Background(), userID, comicID, req)

	assert.False(t, result.Success)
	assert.Equal(t, http.StatusInternalServerError, result.HttpStatus)
	repo.AssertExpectations(t)
}

// ---------------------------------------------------------------------------
// ListRatings
// ---------------------------------------------------------------------------

func TestListRatings_Success(t *testing.T) {
	repo := new(MockRatingRepo)
	svc := newTestService(repo)

	userID := uuid.New()
	ratings := []*model.Rating{sampleRating(userID, uuid.New()), sampleRating(userID, uuid.New())}
	paging := &common.Paging{Page: 1, Limit: 20}

	repo.On("FindPaginated", mock.Anything, mock.Anything, paging, mock.Anything).Return(ratings, int64(2), nil)

	result := svc.ListRatings(context.Background(), userID, paging)

	assert.True(t, result.Success)
	assert.Equal(t, http.StatusOK, result.HttpStatus)
	repo.AssertExpectations(t)
}

func TestListRatings_DBError(t *testing.T) {
	repo := new(MockRatingRepo)
	svc := newTestService(repo)

	paging := &common.Paging{Page: 1, Limit: 20}
	repo.On("FindPaginated", mock.Anything, mock.Anything, paging, mock.Anything).Return(nil, int64(0), errors.New("db error"))

	result := svc.ListRatings(context.Background(), uuid.New(), paging)

	assert.False(t, result.Success)
	assert.Equal(t, http.StatusInternalServerError, result.HttpStatus)
	repo.AssertExpectations(t)
}

// ---------------------------------------------------------------------------
// UpdateRating
// ---------------------------------------------------------------------------

func TestUpdateRating_Success(t *testing.T) {
	repo := new(MockRatingRepo)
	svc := newTestService(repo)

	userID := uuid.New()
	rating := sampleRating(userID, uuid.New())
	req := &ratingrequest.UpdateRatingRequest{Score: 5}

	repo.On("FindOne", mock.Anything, mock.Anything, mock.Anything).Return(rating, nil)
	repo.On("Update", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	result := svc.UpdateRating(context.Background(), userID, rating.ID, req)

	assert.True(t, result.Success)
	assert.Equal(t, http.StatusOK, result.HttpStatus)
	assert.Equal(t, "Rating updated successfully", result.Message)
	repo.AssertExpectations(t)
}

func TestUpdateRating_NotFound(t *testing.T) {
	repo := new(MockRatingRepo)
	svc := newTestService(repo)

	req := &ratingrequest.UpdateRatingRequest{Score: 5}
	repo.On("FindOne", mock.Anything, mock.Anything, mock.Anything).Return(nil, gorm.ErrRecordNotFound)

	result := svc.UpdateRating(context.Background(), uuid.New(), uuid.New(), req)

	assert.False(t, result.Success)
	assert.Equal(t, http.StatusBadRequest, result.HttpStatus)
	assert.Equal(t, "Rating not found", result.Message)
	repo.AssertExpectations(t)
}

func TestUpdateRating_UpdateDBError(t *testing.T) {
	repo := new(MockRatingRepo)
	svc := newTestService(repo)

	userID := uuid.New()
	rating := sampleRating(userID, uuid.New())
	req := &ratingrequest.UpdateRatingRequest{Score: 5}

	repo.On("FindOne", mock.Anything, mock.Anything, mock.Anything).Return(rating, nil)
	repo.On("Update", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("update failed"))

	result := svc.UpdateRating(context.Background(), userID, rating.ID, req)

	assert.False(t, result.Success)
	assert.Equal(t, http.StatusInternalServerError, result.HttpStatus)
	repo.AssertExpectations(t)
}

// ---------------------------------------------------------------------------
// DeleteRating
// ---------------------------------------------------------------------------

func TestDeleteRating_Success(t *testing.T) {
	repo := new(MockRatingRepo)
	svc := newTestService(repo)

	userID := uuid.New()
	rating := sampleRating(userID, uuid.New())

	repo.On("FindOne", mock.Anything, mock.Anything, mock.Anything).Return(rating, nil)
	repo.On("DeleteSoft", mock.Anything, mock.Anything).Return(nil)

	result := svc.DeleteRating(context.Background(), userID, rating.ID)

	assert.True(t, result.Success)
	assert.Equal(t, http.StatusOK, result.HttpStatus)
	assert.Equal(t, "Rating deleted successfully", result.Message)
	repo.AssertExpectations(t)
}

func TestDeleteRating_NotFound(t *testing.T) {
	repo := new(MockRatingRepo)
	svc := newTestService(repo)

	repo.On("FindOne", mock.Anything, mock.Anything, mock.Anything).Return(nil, gorm.ErrRecordNotFound)

	result := svc.DeleteRating(context.Background(), uuid.New(), uuid.New())

	assert.False(t, result.Success)
	assert.Equal(t, http.StatusBadRequest, result.HttpStatus)
	assert.Equal(t, "Rating not found", result.Message)
	repo.AssertExpectations(t)
}

func TestDeleteRating_DeleteDBError(t *testing.T) {
	repo := new(MockRatingRepo)
	svc := newTestService(repo)

	userID := uuid.New()
	rating := sampleRating(userID, uuid.New())

	repo.On("FindOne", mock.Anything, mock.Anything, mock.Anything).Return(rating, nil)
	repo.On("DeleteSoft", mock.Anything, mock.Anything).Return(errors.New("delete failed"))

	result := svc.DeleteRating(context.Background(), userID, rating.ID)

	assert.False(t, result.Success)
	assert.Equal(t, http.StatusInternalServerError, result.HttpStatus)
	repo.AssertExpectations(t)
}

// ---------------------------------------------------------------------------
// GetAverageRating
// ---------------------------------------------------------------------------

func TestGetAverageRating_Success(t *testing.T) {
	repo := new(MockRatingRepo)
	svc := newTestService(repo)

	comicID := uuid.New()
	repo.On("GetAverageRatingByComicID", mock.Anything, comicID).Return(4.5, int64(10), nil)

	result := svc.GetAverageRating(context.Background(), comicID)

	assert.True(t, result.Success)
	assert.Equal(t, http.StatusOK, result.HttpStatus)
	assert.Equal(t, "Average rating retrieved successfully", result.Message)
	repo.AssertExpectations(t)
}

func TestGetAverageRating_DBError(t *testing.T) {
	repo := new(MockRatingRepo)
	svc := newTestService(repo)

	comicID := uuid.New()
	repo.On("GetAverageRatingByComicID", mock.Anything, comicID).Return(0.0, int64(0), errors.New("db error"))

	result := svc.GetAverageRating(context.Background(), comicID)

	assert.False(t, result.Success)
	assert.Equal(t, http.StatusInternalServerError, result.HttpStatus)
	repo.AssertExpectations(t)
}
