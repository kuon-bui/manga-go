package readinghistoryservice_test

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"manga-go/internal/pkg/common"
	"manga-go/internal/pkg/logger"
	"manga-go/internal/pkg/model"
	readinghistoryrequest "manga-go/internal/pkg/request/reading_history"
	readinghistoryservice "manga-go/internal/pkg/services/reading_history"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// MockReadingHistoryRepo implements readinghistoryservice.ReadingHistoryRepository.
type MockReadingHistoryRepo struct {
	mock.Mock
}

func (m *MockReadingHistoryRepo) Create(ctx context.Context, rh *model.ReadingHistory) error {
	args := m.Called(ctx, rh)
	return args.Error(0)
}

func (m *MockReadingHistoryRepo) FindOne(ctx context.Context, conditions []any, moreKeys map[string]common.MoreKeyOption) (*model.ReadingHistory, error) {
	args := m.Called(ctx, conditions, moreKeys)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.ReadingHistory), args.Error(1)
}

func (m *MockReadingHistoryRepo) Update(ctx context.Context, conditions []any, data map[string]any) error {
	args := m.Called(ctx, conditions, data)
	return args.Error(0)
}

func (m *MockReadingHistoryRepo) DeleteSoft(ctx context.Context, conditions []any) error {
	args := m.Called(ctx, conditions)
	return args.Error(0)
}

func (m *MockReadingHistoryRepo) FindPaginated(ctx context.Context, conditions []any, paging *common.Paging, moreKeys map[string]common.MoreKeyOption) ([]*model.ReadingHistory, int64, error) {
	args := m.Called(ctx, conditions, paging, moreKeys)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*model.ReadingHistory), args.Get(1).(int64), args.Error(2)
}

func newTestService(repo *MockReadingHistoryRepo) *readinghistoryservice.ReadingHistoryService {
	return readinghistoryservice.NewReadingHistoryServiceWithRepo(logger.NewLogger(), repo)
}

func sampleReadingHistory(userID uuid.UUID) *model.ReadingHistory {
	now := time.Now()
	return &model.ReadingHistory{
		SqlModel: common.SqlModel{
			ID:        uuid.New(),
			CreatedAt: &now,
			UpdatedAt: &now,
		},
		UserID:     userID,
		ChapterID:  uuid.New(),
		ComicID:    uuid.New(),
		LastReadAt: &now,
	}
}

// ---------------------------------------------------------------------------
// CreateReadingHistory
// ---------------------------------------------------------------------------

func TestCreateReadingHistory_Success(t *testing.T) {
	repo := new(MockReadingHistoryRepo)
	svc := newTestService(repo)

	userID := uuid.New()
	req := &readinghistoryrequest.CreateReadingHistoryRequest{
		ChapterID: uuid.New(),
		ComicID:   uuid.New(),
	}
	repo.On("Create", mock.Anything, mock.AnythingOfType("*model.ReadingHistory")).Return(nil)

	result := svc.CreateReadingHistory(context.Background(), userID, req)

	assert.True(t, result.Success)
	assert.Equal(t, http.StatusOK, result.HttpStatus)
	assert.Equal(t, "Reading history created successfully", result.Message)
	repo.AssertExpectations(t)
}

func TestCreateReadingHistory_DBError(t *testing.T) {
	repo := new(MockReadingHistoryRepo)
	svc := newTestService(repo)

	userID := uuid.New()
	req := &readinghistoryrequest.CreateReadingHistoryRequest{
		ChapterID: uuid.New(),
		ComicID:   uuid.New(),
	}
	repo.On("Create", mock.Anything, mock.AnythingOfType("*model.ReadingHistory")).Return(errors.New("db error"))

	result := svc.CreateReadingHistory(context.Background(), userID, req)

	assert.False(t, result.Success)
	assert.Equal(t, http.StatusInternalServerError, result.HttpStatus)
	repo.AssertExpectations(t)
}

// ---------------------------------------------------------------------------
// GetReadingHistory
// ---------------------------------------------------------------------------

func TestGetReadingHistory_Success(t *testing.T) {
	repo := new(MockReadingHistoryRepo)
	svc := newTestService(repo)

	rh := sampleReadingHistory(uuid.New())
	repo.On("FindOne", mock.Anything, mock.Anything, mock.Anything).Return(rh, nil)

	result := svc.GetReadingHistory(context.Background(), rh.ID)

	assert.True(t, result.Success)
	assert.Equal(t, http.StatusOK, result.HttpStatus)
	assert.Equal(t, "Reading history retrieved successfully", result.Message)
	repo.AssertExpectations(t)
}

func TestGetReadingHistory_NotFound(t *testing.T) {
	repo := new(MockReadingHistoryRepo)
	svc := newTestService(repo)

	repo.On("FindOne", mock.Anything, mock.Anything, mock.Anything).Return(nil, gorm.ErrRecordNotFound)

	result := svc.GetReadingHistory(context.Background(), uuid.New())

	assert.False(t, result.Success)
	assert.Equal(t, http.StatusBadRequest, result.HttpStatus)
	assert.Equal(t, "ReadingHistory not found", result.Message)
	repo.AssertExpectations(t)
}

func TestGetReadingHistory_DBError(t *testing.T) {
	repo := new(MockReadingHistoryRepo)
	svc := newTestService(repo)

	repo.On("FindOne", mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("timeout"))

	result := svc.GetReadingHistory(context.Background(), uuid.New())

	assert.False(t, result.Success)
	assert.Equal(t, http.StatusInternalServerError, result.HttpStatus)
	repo.AssertExpectations(t)
}

// ---------------------------------------------------------------------------
// ListReadingHistories
// ---------------------------------------------------------------------------

func TestListReadingHistories_Success(t *testing.T) {
	repo := new(MockReadingHistoryRepo)
	svc := newTestService(repo)

	userID := uuid.New()
	histories := []*model.ReadingHistory{sampleReadingHistory(userID), sampleReadingHistory(userID)}
	paging := &common.Paging{Page: 1, Limit: 20}
	repo.On("FindPaginated", mock.Anything, mock.Anything, paging, mock.Anything).Return(histories, int64(2), nil)

	result := svc.ListReadingHistories(context.Background(), userID, paging)

	assert.True(t, result.Success)
	assert.Equal(t, http.StatusOK, result.HttpStatus)
	repo.AssertExpectations(t)
}

func TestListReadingHistories_DBError(t *testing.T) {
	repo := new(MockReadingHistoryRepo)
	svc := newTestService(repo)

	paging := &common.Paging{Page: 1, Limit: 20}
	repo.On("FindPaginated", mock.Anything, mock.Anything, paging, mock.Anything).Return(nil, int64(0), errors.New("db error"))

	result := svc.ListReadingHistories(context.Background(), uuid.New(), paging)

	assert.False(t, result.Success)
	assert.Equal(t, http.StatusInternalServerError, result.HttpStatus)
	repo.AssertExpectations(t)
}

// ---------------------------------------------------------------------------
// UpdateReadingHistory
// ---------------------------------------------------------------------------

func TestUpdateReadingHistory_Success(t *testing.T) {
	repo := new(MockReadingHistoryRepo)
	svc := newTestService(repo)

	rh := sampleReadingHistory(uuid.New())
	now := time.Now()
	req := &readinghistoryrequest.UpdateReadingHistoryRequest{LastReadAt: &now}

	repo.On("FindOne", mock.Anything, mock.Anything, mock.Anything).Return(rh, nil)
	repo.On("Update", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	result := svc.UpdateReadingHistory(context.Background(), rh.ID, req)

	assert.True(t, result.Success)
	assert.Equal(t, http.StatusOK, result.HttpStatus)
	assert.Equal(t, "Reading history updated successfully", result.Message)
	repo.AssertExpectations(t)
}

func TestUpdateReadingHistory_UpdatesNowWhenNoTimestamp(t *testing.T) {
	repo := new(MockReadingHistoryRepo)
	svc := newTestService(repo)

	rh := sampleReadingHistory(uuid.New())
	req := &readinghistoryrequest.UpdateReadingHistoryRequest{LastReadAt: nil}

	repo.On("FindOne", mock.Anything, mock.Anything, mock.Anything).Return(rh, nil)
	repo.On("Update", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	result := svc.UpdateReadingHistory(context.Background(), rh.ID, req)

	assert.True(t, result.Success)
	assert.Equal(t, http.StatusOK, result.HttpStatus)
	repo.AssertExpectations(t)
}

func TestUpdateReadingHistory_NotFound(t *testing.T) {
	repo := new(MockReadingHistoryRepo)
	svc := newTestService(repo)

	repo.On("FindOne", mock.Anything, mock.Anything, mock.Anything).Return(nil, gorm.ErrRecordNotFound)

	result := svc.UpdateReadingHistory(context.Background(), uuid.New(), &readinghistoryrequest.UpdateReadingHistoryRequest{})

	assert.False(t, result.Success)
	assert.Equal(t, http.StatusBadRequest, result.HttpStatus)
	assert.Equal(t, "ReadingHistory not found", result.Message)
	repo.AssertExpectations(t)
}

func TestUpdateReadingHistory_UpdateDBError(t *testing.T) {
	repo := new(MockReadingHistoryRepo)
	svc := newTestService(repo)

	rh := sampleReadingHistory(uuid.New())
	repo.On("FindOne", mock.Anything, mock.Anything, mock.Anything).Return(rh, nil)
	repo.On("Update", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("update failed"))

	result := svc.UpdateReadingHistory(context.Background(), rh.ID, &readinghistoryrequest.UpdateReadingHistoryRequest{})

	assert.False(t, result.Success)
	assert.Equal(t, http.StatusInternalServerError, result.HttpStatus)
	repo.AssertExpectations(t)
}

// ---------------------------------------------------------------------------
// DeleteReadingHistory
// ---------------------------------------------------------------------------

func TestDeleteReadingHistory_Success(t *testing.T) {
	repo := new(MockReadingHistoryRepo)
	svc := newTestService(repo)

	rh := sampleReadingHistory(uuid.New())
	repo.On("FindOne", mock.Anything, mock.Anything, mock.Anything).Return(rh, nil)
	repo.On("DeleteSoft", mock.Anything, mock.Anything).Return(nil)

	result := svc.DeleteReadingHistory(context.Background(), rh.ID)

	assert.True(t, result.Success)
	assert.Equal(t, http.StatusOK, result.HttpStatus)
	assert.Equal(t, "Reading history deleted successfully", result.Message)
	repo.AssertExpectations(t)
}

func TestDeleteReadingHistory_NotFound(t *testing.T) {
	repo := new(MockReadingHistoryRepo)
	svc := newTestService(repo)

	repo.On("FindOne", mock.Anything, mock.Anything, mock.Anything).Return(nil, gorm.ErrRecordNotFound)

	result := svc.DeleteReadingHistory(context.Background(), uuid.New())

	assert.False(t, result.Success)
	assert.Equal(t, http.StatusBadRequest, result.HttpStatus)
	assert.Equal(t, "ReadingHistory not found", result.Message)
	repo.AssertExpectations(t)
}

func TestDeleteReadingHistory_DeleteDBError(t *testing.T) {
	repo := new(MockReadingHistoryRepo)
	svc := newTestService(repo)

	rh := sampleReadingHistory(uuid.New())
	repo.On("FindOne", mock.Anything, mock.Anything, mock.Anything).Return(rh, nil)
	repo.On("DeleteSoft", mock.Anything, mock.Anything).Return(errors.New("delete failed"))

	result := svc.DeleteReadingHistory(context.Background(), rh.ID)

	assert.False(t, result.Success)
	assert.Equal(t, http.StatusInternalServerError, result.HttpStatus)
	repo.AssertExpectations(t)
}
