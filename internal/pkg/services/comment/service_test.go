package commentservice_test

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"manga-go/internal/pkg/common"
	"manga-go/internal/pkg/logger"
	"manga-go/internal/pkg/model"
	commentrequest "manga-go/internal/pkg/request/comment"
	commentservice "manga-go/internal/pkg/services/comment"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// ---------------------------------------------------------------------------
// Mock implementations
// ---------------------------------------------------------------------------

type MockCommentRepo struct{ mock.Mock }

func (m *MockCommentRepo) Create(ctx context.Context, comment *model.Comment) error {
	args := m.Called(ctx, comment)
	return args.Error(0)
}

func (m *MockCommentRepo) FindOne(ctx context.Context, conditions []any, moreKeys map[string]common.MoreKeyOption) (*model.Comment, error) {
	args := m.Called(ctx, conditions, moreKeys)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Comment), args.Error(1)
}

func (m *MockCommentRepo) Update(ctx context.Context, conditions []any, data map[string]any) error {
	args := m.Called(ctx, conditions, data)
	return args.Error(0)
}

func (m *MockCommentRepo) DeleteSoft(ctx context.Context, conditions []any) error {
	args := m.Called(ctx, conditions)
	return args.Error(0)
}

func (m *MockCommentRepo) FindPaginated(ctx context.Context, conditions []any, paging *common.Paging, moreKeys map[string]common.MoreKeyOption) ([]*model.Comment, int64, error) {
	args := m.Called(ctx, conditions, paging, moreKeys)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*model.Comment), args.Get(1).(int64), args.Error(2)
}

type MockChapterRepo struct{ mock.Mock }

func (m *MockChapterRepo) FindOne(ctx context.Context, conditions []any, moreKeys map[string]common.MoreKeyOption) (*model.Chapter, error) {
	args := m.Called(ctx, conditions, moreKeys)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Chapter), args.Error(1)
}

type MockReactionRepo struct{ mock.Mock }

func (m *MockReactionRepo) Create(ctx context.Context, reaction *model.Reaction) error {
	args := m.Called(ctx, reaction)
	return args.Error(0)
}

func (m *MockReactionRepo) DeleteSoft(ctx context.Context, conditions []any) error {
	args := m.Called(ctx, conditions)
	return args.Error(0)
}

func (m *MockReactionRepo) ExistsByCommentIdAndUserId(ctx context.Context, commentId, userId uuid.UUID) (bool, error) {
	args := m.Called(ctx, commentId, userId)
	return args.Bool(0), args.Error(1)
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

func newTestService(commentRepo *MockCommentRepo, chapterRepo *MockChapterRepo, reactionRepo *MockReactionRepo) *commentservice.CommentService {
	return commentservice.NewCommentServiceWithRepos(logger.NewLogger(), commentRepo, chapterRepo, reactionRepo)
}

func sampleChapter(comicID uuid.UUID) *model.Chapter {
	now := time.Now()
	return &model.Chapter{
		SqlModel: common.SqlModel{
			ID:        uuid.New(),
			CreatedAt: &now,
			UpdatedAt: &now,
		},
		ComicID: comicID,
		Number:  "1",
		Title:   "First Chapter",
		Slug:    "chapter-1",
	}
}

func sampleComment(userID, chapterID, comicID uuid.UUID) *model.Comment {
	now := time.Now()
	return &model.Comment{
		SqlModel: common.SqlModel{
			ID:        uuid.New(),
			CreatedAt: &now,
			UpdatedAt: &now,
		},
		UserId:    userID,
		ChapterId: chapterID,
		ComicId:   comicID,
		Content:   "Great chapter!",
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
// CreateComment
// ---------------------------------------------------------------------------

func TestCreateComment_Success(t *testing.T) {
	commentRepo := new(MockCommentRepo)
	chapterRepo := new(MockChapterRepo)
	svc := newTestService(commentRepo, chapterRepo, new(MockReactionRepo))

	comicID := uuid.New()
	chapter := sampleChapter(comicID)
	userID := uuid.New()
	req := &commentrequest.CreateCommentRequest{ChapterID: chapter.ID, Content: "Awesome!"}

	chapterRepo.On("FindOne", mock.Anything, mock.Anything, mock.Anything).Return(chapter, nil)
	commentRepo.On("Create", mock.Anything, mock.AnythingOfType("*model.Comment")).Return(nil)

	result := svc.CreateComment(context.Background(), userID, req)

	assert.True(t, result.Success)
	assert.Equal(t, http.StatusOK, result.HttpStatus)
	assert.Equal(t, "Comment created successfully", result.Message)
	chapterRepo.AssertExpectations(t)
	commentRepo.AssertExpectations(t)
}

func TestCreateComment_ChapterNotFound(t *testing.T) {
	commentRepo := new(MockCommentRepo)
	chapterRepo := new(MockChapterRepo)
	svc := newTestService(commentRepo, chapterRepo, new(MockReactionRepo))

	chapterRepo.On("FindOne", mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("not found"))

	result := svc.CreateComment(context.Background(), uuid.New(), &commentrequest.CreateCommentRequest{ChapterID: uuid.New(), Content: "Hi"})

	assert.False(t, result.Success)
	assert.Equal(t, http.StatusInternalServerError, result.HttpStatus)
	chapterRepo.AssertExpectations(t)
}

func TestCreateComment_DBError(t *testing.T) {
	commentRepo := new(MockCommentRepo)
	chapterRepo := new(MockChapterRepo)
	svc := newTestService(commentRepo, chapterRepo, new(MockReactionRepo))

	comicID := uuid.New()
	chapter := sampleChapter(comicID)
	req := &commentrequest.CreateCommentRequest{ChapterID: chapter.ID, Content: "Awesome!"}

	chapterRepo.On("FindOne", mock.Anything, mock.Anything, mock.Anything).Return(chapter, nil)
	commentRepo.On("Create", mock.Anything, mock.AnythingOfType("*model.Comment")).Return(errors.New("db error"))

	result := svc.CreateComment(context.Background(), uuid.New(), req)

	assert.False(t, result.Success)
	assert.Equal(t, http.StatusInternalServerError, result.HttpStatus)
	commentRepo.AssertExpectations(t)
}

// ---------------------------------------------------------------------------
// GetComment
// ---------------------------------------------------------------------------

func TestGetComment_Success(t *testing.T) {
	commentRepo := new(MockCommentRepo)
	svc := newTestService(commentRepo, new(MockChapterRepo), new(MockReactionRepo))

	comment := sampleComment(uuid.New(), uuid.New(), uuid.New())
	commentRepo.On("FindOne", mock.Anything, mock.Anything, mock.Anything).Return(comment, nil)

	result := svc.GetComment(context.Background(), comment.ID)

	assert.True(t, result.Success)
	assert.Equal(t, http.StatusOK, result.HttpStatus)
	assert.Equal(t, "Comment retrieved successfully", result.Message)
	commentRepo.AssertExpectations(t)
}

func TestGetComment_NotFound(t *testing.T) {
	commentRepo := new(MockCommentRepo)
	svc := newTestService(commentRepo, new(MockChapterRepo), new(MockReactionRepo))

	commentRepo.On("FindOne", mock.Anything, mock.Anything, mock.Anything).Return(nil, gorm.ErrRecordNotFound)

	result := svc.GetComment(context.Background(), uuid.New())

	assert.False(t, result.Success)
	assert.Equal(t, http.StatusBadRequest, result.HttpStatus)
	assert.Equal(t, "Comment not found", result.Message)
	commentRepo.AssertExpectations(t)
}

// ---------------------------------------------------------------------------
// ListComments
// ---------------------------------------------------------------------------

func TestListComments_Success(t *testing.T) {
	commentRepo := new(MockCommentRepo)
	svc := newTestService(commentRepo, new(MockChapterRepo), new(MockReactionRepo))

	chapterID := uuid.New()
	comicID := uuid.New()
	comments := []*model.Comment{
		sampleComment(uuid.New(), chapterID, comicID),
		sampleComment(uuid.New(), chapterID, comicID),
	}
	req := &commentrequest.ListCommentsRequest{
		Paging:    common.Paging{Page: 1, Limit: 20},
		ChapterId: chapterID,
	}
	commentRepo.On("FindPaginated", mock.Anything, mock.Anything, &req.Paging, mock.Anything).Return(comments, int64(2), nil)

	result := svc.ListComments(context.Background(), req)

	assert.True(t, result.Success)
	assert.Equal(t, http.StatusOK, result.HttpStatus)
	commentRepo.AssertExpectations(t)
}

func TestListComments_DBError(t *testing.T) {
	commentRepo := new(MockCommentRepo)
	svc := newTestService(commentRepo, new(MockChapterRepo), new(MockReactionRepo))

	req := &commentrequest.ListCommentsRequest{
		Paging:    common.Paging{Page: 1, Limit: 20},
		ChapterId: uuid.New(),
	}
	commentRepo.On("FindPaginated", mock.Anything, mock.Anything, &req.Paging, mock.Anything).Return(nil, int64(0), errors.New("db error"))

	result := svc.ListComments(context.Background(), req)

	assert.False(t, result.Success)
	assert.Equal(t, http.StatusInternalServerError, result.HttpStatus)
	commentRepo.AssertExpectations(t)
}

// ---------------------------------------------------------------------------
// UpdateComment
// ---------------------------------------------------------------------------

func TestUpdateComment_Success(t *testing.T) {
	commentRepo := new(MockCommentRepo)
	svc := newTestService(commentRepo, new(MockChapterRepo), new(MockReactionRepo))

	comment := sampleComment(uuid.New(), uuid.New(), uuid.New())
	req := &commentrequest.UpdateCommentRequest{Content: "Updated content"}

	commentRepo.On("FindOne", mock.Anything, mock.Anything, mock.Anything).Return(comment, nil)
	commentRepo.On("Update", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	result := svc.UpdateComment(context.Background(), comment.ID, req)

	assert.True(t, result.Success)
	assert.Equal(t, http.StatusOK, result.HttpStatus)
	assert.Equal(t, "Comment updated successfully", result.Message)
	commentRepo.AssertExpectations(t)
}

func TestUpdateComment_NotFound(t *testing.T) {
	commentRepo := new(MockCommentRepo)
	svc := newTestService(commentRepo, new(MockChapterRepo), new(MockReactionRepo))

	req := &commentrequest.UpdateCommentRequest{Content: "Updated"}
	commentRepo.On("FindOne", mock.Anything, mock.Anything, mock.Anything).Return(nil, gorm.ErrRecordNotFound)

	result := svc.UpdateComment(context.Background(), uuid.New(), req)

	assert.False(t, result.Success)
	assert.Equal(t, http.StatusBadRequest, result.HttpStatus)
	assert.Equal(t, "Comment not found", result.Message)
	commentRepo.AssertExpectations(t)
}

func TestUpdateComment_UpdateDBError(t *testing.T) {
	commentRepo := new(MockCommentRepo)
	svc := newTestService(commentRepo, new(MockChapterRepo), new(MockReactionRepo))

	comment := sampleComment(uuid.New(), uuid.New(), uuid.New())
	req := &commentrequest.UpdateCommentRequest{Content: "Updated content"}

	commentRepo.On("FindOne", mock.Anything, mock.Anything, mock.Anything).Return(comment, nil)
	commentRepo.On("Update", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("update failed"))

	result := svc.UpdateComment(context.Background(), comment.ID, req)

	assert.False(t, result.Success)
	assert.Equal(t, http.StatusInternalServerError, result.HttpStatus)
	commentRepo.AssertExpectations(t)
}

// ---------------------------------------------------------------------------
// DeleteComment
// ---------------------------------------------------------------------------

func TestDeleteComment_Success(t *testing.T) {
	commentRepo := new(MockCommentRepo)
	svc := newTestService(commentRepo, new(MockChapterRepo), new(MockReactionRepo))

	comment := sampleComment(uuid.New(), uuid.New(), uuid.New())
	commentRepo.On("FindOne", mock.Anything, mock.Anything, mock.Anything).Return(comment, nil)
	commentRepo.On("DeleteSoft", mock.Anything, mock.Anything).Return(nil)

	result := svc.DeleteComment(context.Background(), comment.ID)

	assert.True(t, result.Success)
	assert.Equal(t, http.StatusOK, result.HttpStatus)
	assert.Equal(t, "Comment deleted successfully", result.Message)
	commentRepo.AssertExpectations(t)
}

func TestDeleteComment_NotFound(t *testing.T) {
	commentRepo := new(MockCommentRepo)
	svc := newTestService(commentRepo, new(MockChapterRepo), new(MockReactionRepo))

	commentRepo.On("FindOne", mock.Anything, mock.Anything, mock.Anything).Return(nil, gorm.ErrRecordNotFound)

	result := svc.DeleteComment(context.Background(), uuid.New())

	assert.False(t, result.Success)
	assert.Equal(t, http.StatusBadRequest, result.HttpStatus)
	assert.Equal(t, "Comment not found", result.Message)
	commentRepo.AssertExpectations(t)
}

func TestDeleteComment_DeleteDBError(t *testing.T) {
	commentRepo := new(MockCommentRepo)
	svc := newTestService(commentRepo, new(MockChapterRepo), new(MockReactionRepo))

	comment := sampleComment(uuid.New(), uuid.New(), uuid.New())
	commentRepo.On("FindOne", mock.Anything, mock.Anything, mock.Anything).Return(comment, nil)
	commentRepo.On("DeleteSoft", mock.Anything, mock.Anything).Return(errors.New("delete failed"))

	result := svc.DeleteComment(context.Background(), comment.ID)

	assert.False(t, result.Success)
	assert.Equal(t, http.StatusInternalServerError, result.HttpStatus)
	commentRepo.AssertExpectations(t)
}

// ---------------------------------------------------------------------------
// HandleReaction — add reaction (user has not reacted yet)
// ---------------------------------------------------------------------------

func TestHandleReaction_AddReaction_Success(t *testing.T) {
	commentRepo := new(MockCommentRepo)
	reactionRepo := new(MockReactionRepo)
	svc := newTestService(commentRepo, new(MockChapterRepo), reactionRepo)

	user := sampleUser()
	comment := sampleComment(user.ID, uuid.New(), uuid.New())
	req := &commentrequest.AddReactionRequest{Type: "like"}

	commentRepo.On("FindOne", mock.Anything, mock.Anything, mock.Anything).Return(comment, nil)
	reactionRepo.On("ExistsByCommentIdAndUserId", mock.Anything, comment.ID, user.ID).Return(false, nil)
	reactionRepo.On("Create", mock.Anything, mock.AnythingOfType("*model.Reaction")).Return(nil)

	result := svc.HandleReaction(context.Background(), user, comment.ID, req)

	assert.True(t, result.Success)
	assert.Equal(t, http.StatusOK, result.HttpStatus)
	assert.Equal(t, "Reaction added successfully", result.Message)
	commentRepo.AssertExpectations(t)
	reactionRepo.AssertExpectations(t)
}

// HandleReaction — remove reaction (user already reacted)
func TestHandleReaction_RemoveReaction_Success(t *testing.T) {
	commentRepo := new(MockCommentRepo)
	reactionRepo := new(MockReactionRepo)
	svc := newTestService(commentRepo, new(MockChapterRepo), reactionRepo)

	user := sampleUser()
	comment := sampleComment(user.ID, uuid.New(), uuid.New())
	req := &commentrequest.AddReactionRequest{Type: "like"}

	commentRepo.On("FindOne", mock.Anything, mock.Anything, mock.Anything).Return(comment, nil)
	reactionRepo.On("ExistsByCommentIdAndUserId", mock.Anything, comment.ID, user.ID).Return(true, nil)
	reactionRepo.On("DeleteSoft", mock.Anything, mock.Anything).Return(nil)

	result := svc.HandleReaction(context.Background(), user, comment.ID, req)

	assert.True(t, result.Success)
	assert.Equal(t, http.StatusOK, result.HttpStatus)
	assert.Equal(t, "Reaction removed successfully", result.Message)
	commentRepo.AssertExpectations(t)
	reactionRepo.AssertExpectations(t)
}

func TestHandleReaction_CommentNotFound(t *testing.T) {
	commentRepo := new(MockCommentRepo)
	svc := newTestService(commentRepo, new(MockChapterRepo), new(MockReactionRepo))

	user := sampleUser()
	commentRepo.On("FindOne", mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("db error"))

	result := svc.HandleReaction(context.Background(), user, uuid.New(), &commentrequest.AddReactionRequest{Type: "like"})

	assert.False(t, result.Success)
	assert.Equal(t, http.StatusInternalServerError, result.HttpStatus)
	commentRepo.AssertExpectations(t)
}
