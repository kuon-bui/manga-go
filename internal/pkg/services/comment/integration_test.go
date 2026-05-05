//go:build integration

package commentservice

import (
	"context"
	"net/http"
	"reflect"
	"testing"

	"manga-go/internal/pkg/logger"
	"manga-go/internal/pkg/model"
	chapterrepo "manga-go/internal/pkg/repo/chapter"
	commentrepo "manga-go/internal/pkg/repo/comment"
	reactionrepo "manga-go/internal/pkg/repo/reaction"
	commentrequest "manga-go/internal/pkg/request/comment"
	"manga-go/internal/pkg/testutil"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func newCommentServiceIntegration(t *testing.T) (*CommentService, *gorm.DB, uuid.UUID) {
	t.Helper()

	db := testutil.NewSQLiteDB(t)
	tx := db.Begin()
	if tx.Error != nil {
		t.Fatalf("failed to begin transaction: %v", tx.Error)
	}
	t.Cleanup(func() {
		_ = tx.Rollback().Error
	})

	testutil.MustSyncSchemas(t, tx,
		&testutil.Chapter{},
		&testutil.User{},
		&testutil.Comment{},
		&testutil.CommentReaction{},
	)

	chapterID := uuid.New()
	comicID := uuid.New()
	if err := tx.Exec(`INSERT INTO chapters (id, comic_id) VALUES (?, ?)`, chapterID, comicID).Error; err != nil {
		t.Fatalf("failed to seed chapter: %v", err)
	}

	s := &CommentService{
		logger:       logger.NewLogger(),
		commentRepo:  commentrepo.NewCommentRepo(tx),
		chapterRepo:  chapterrepo.NewChapterRepo(tx, nil),
		reactionRepo: reactionrepo.NewReactionRepo(tx),
	}

	return s, tx, chapterID
}

func commentPaginationTotalFromData(data any) int64 {
	v := reflect.ValueOf(data)
	if !v.IsValid() {
		return -1
	}

	field := v.FieldByName("Total")
	if !field.IsValid() || field.Kind() != reflect.Int64 {
		return -1
	}

	return field.Int()
}

func TestCommentServiceIntegrationFullFlow(t *testing.T) {
	s, db, chapterID := newCommentServiceIntegration(t)
	ctx := context.Background()
	pageIndex := 3

	createRes := s.CreateComment(ctx, uuid.New(), &commentrequest.CreateCommentRequest{
		ChapterID: &chapterID,
		Content:   "first comment",
		PageIndex: &pageIndex,
	})
	if !createRes.Success {
		t.Fatalf("expected create success, got: %s", createRes.Message)
	}

	listRes := s.ListComments(ctx, &commentrequest.ListCommentsRequest{ChapterId: chapterID.String()})
	if !listRes.Success {
		t.Fatalf("expected list success, got: %s", listRes.Message)
	}
	if total := commentPaginationTotalFromData(listRes.Data); total != 1 {
		t.Fatalf("expected total 1, got %d", total)
	}

	commentID := testutil.MustReadUUID(t, db, "SELECT id FROM comments WHERE chapter_id = ? AND deleted_at IS NULL", chapterID)
	if commentID == uuid.Nil {
		t.Fatalf("expected persisted comment id")
	}

	getRes := s.GetComment(ctx, commentID)
	if !getRes.Success {
		t.Fatalf("expected get success, got: %s", getRes.Message)
	}

	newPageIndex := 5
	updateRes := s.UpdateComment(ctx, commentID, &commentrequest.UpdateCommentRequest{
		Content:   "updated comment",
		PageIndex: &newPageIndex,
	})
	if !updateRes.Success {
		t.Fatalf("expected update success, got: %s", updateRes.Message)
	}

	deleteRes := s.DeleteComment(ctx, commentID)
	if !deleteRes.Success {
		t.Fatalf("expected delete success, got: %s", deleteRes.Message)
	}

	notFoundRes := s.GetComment(ctx, commentID)
	if notFoundRes.Success {
		t.Fatalf("expected not found after soft delete")
	}
	if notFoundRes.Message != "Comment not found" {
		t.Fatalf("unexpected message: %s", notFoundRes.Message)
	}
}

func TestCommentServiceIntegrationCreateCommentChapterNotFound(t *testing.T) {
	s, _, _ := newCommentServiceIntegration(t)
	ctx := context.Background()
	pageIndex := 1

	missingChapterID := uuid.New()
	res := s.CreateComment(ctx, uuid.New(), &commentrequest.CreateCommentRequest{
		ChapterID: &missingChapterID,
		Content:   "comment for missing chapter",
		PageIndex: &pageIndex,
	})

	if res.Success {
		t.Fatalf("expected create to fail when chapter does not exist")
	}
	if res.HttpStatus != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", res.HttpStatus)
	}
	if res.Message != "Chapter not found" {
		t.Fatalf("unexpected message: %s", res.Message)
	}
}

func TestCommentServiceIntegrationHandleReactionToggleFlow(t *testing.T) {
	s, db, chapterID := newCommentServiceIntegration(t)
	ctx := context.Background()
	userID := uuid.New()
	pageIndex := 2

	createRes := s.CreateComment(ctx, userID, &commentrequest.CreateCommentRequest{
		ChapterID: &chapterID,
		Content:   "reaction target",
		PageIndex: &pageIndex,
	})
	if !createRes.Success {
		t.Fatalf("expected create comment success, got: %s", createRes.Message)
	}

	commentID := testutil.MustReadUUID(t, db, "SELECT id FROM comments WHERE chapter_id = ? AND deleted_at IS NULL", chapterID)
	if commentID == uuid.Nil {
		t.Fatalf("expected persisted comment id")
	}

	user := &model.User{}
	user.ID = userID

	addRes := s.HandleReaction(ctx, user, commentID, &commentrequest.AddReactionRequest{Type: "like"})
	if !addRes.Success {
		t.Fatalf("expected add reaction success, got: %s", addRes.Message)
	}
	if addRes.Message != "Reaction added successfully" {
		t.Fatalf("unexpected add message: %s", addRes.Message)
	}

	var activeReactions int64
	if err := db.Raw("SELECT COUNT(*) FROM comment_reactions WHERE comment_id = ? AND user_id = ? AND deleted_at IS NULL", commentID, userID).Scan(&activeReactions).Error; err != nil {
		t.Fatalf("failed to count active reactions: %v", err)
	}
	if activeReactions != 1 {
		t.Fatalf("expected 1 active reaction after add, got %d", activeReactions)
	}

	removeRes := s.HandleReaction(ctx, user, commentID, &commentrequest.AddReactionRequest{Type: "like"})
	if !removeRes.Success {
		t.Fatalf("expected remove reaction success, got: %s", removeRes.Message)
	}
	if removeRes.Message != "Reaction removed successfully" {
		t.Fatalf("unexpected remove message: %s", removeRes.Message)
	}

	if err := db.Raw("SELECT COUNT(*) FROM comment_reactions WHERE comment_id = ? AND user_id = ? AND deleted_at IS NULL", commentID, userID).Scan(&activeReactions).Error; err != nil {
		t.Fatalf("failed to count active reactions after remove: %v", err)
	}
	if activeReactions != 0 {
		t.Fatalf("expected 0 active reaction after remove, got %d", activeReactions)
	}

	reAddRes := s.HandleReaction(ctx, user, commentID, &commentrequest.AddReactionRequest{Type: "love"})
	if !reAddRes.Success {
		t.Fatalf("expected re-add reaction success, got: %s", reAddRes.Message)
	}
	if reAddRes.Message != "Reaction added successfully" {
		t.Fatalf("unexpected re-add message: %s", reAddRes.Message)
	}

	if err := db.Raw("SELECT COUNT(*) FROM comment_reactions WHERE comment_id = ? AND user_id = ? AND deleted_at IS NULL", commentID, userID).Scan(&activeReactions).Error; err != nil {
		t.Fatalf("failed to count active reactions after re-add: %v", err)
	}
	if activeReactions != 1 {
		t.Fatalf("expected 1 active reaction after re-add, got %d", activeReactions)
	}
}

func TestCommentServiceIntegrationHandleReactionCommentNotFound(t *testing.T) {
	s, _, _ := newCommentServiceIntegration(t)
	ctx := context.Background()

	user := &model.User{}
	user.ID = uuid.New()

	res := s.HandleReaction(ctx, user, uuid.New(), &commentrequest.AddReactionRequest{Type: "like"})

	if res.Success {
		t.Fatalf("expected handle reaction to fail when comment does not exist")
	}
	if res.HttpStatus != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", res.HttpStatus)
	}
	if res.Message != "database error" {
		t.Fatalf("unexpected message: %s", res.Message)
	}
}

func TestCommentServiceIntegrationHandleReactionIsUserScoped(t *testing.T) {
	s, db, chapterID := newCommentServiceIntegration(t)
	ctx := context.Background()
	pageIndex := 1
	authorID := uuid.New()
	userOneID := uuid.New()
	userTwoID := uuid.New()

	createRes := s.CreateComment(ctx, authorID, &commentrequest.CreateCommentRequest{
		ChapterID: &chapterID,
		Content:   "scoped reaction target",
		PageIndex: &pageIndex,
	})
	if !createRes.Success {
		t.Fatalf("expected create comment success, got: %s", createRes.Message)
	}

	commentID := testutil.MustReadUUID(t, db, "SELECT id FROM comments WHERE chapter_id = ? AND deleted_at IS NULL", chapterID)

	userOne := &model.User{}
	userOne.ID = userOneID
	userTwo := &model.User{}
	userTwo.ID = userTwoID

	addUserOneRes := s.HandleReaction(ctx, userOne, commentID, &commentrequest.AddReactionRequest{Type: "like"})
	if !addUserOneRes.Success {
		t.Fatalf("expected user one add reaction success, got: %s", addUserOneRes.Message)
	}

	addUserTwoRes := s.HandleReaction(ctx, userTwo, commentID, &commentrequest.AddReactionRequest{Type: "like"})
	if !addUserTwoRes.Success {
		t.Fatalf("expected user two add reaction success, got: %s", addUserTwoRes.Message)
	}

	var activeTotal int64
	if err := db.Raw("SELECT COUNT(*) FROM comment_reactions WHERE comment_id = ? AND deleted_at IS NULL", commentID).Scan(&activeTotal).Error; err != nil {
		t.Fatalf("failed to count active reactions: %v", err)
	}
	if activeTotal != 2 {
		t.Fatalf("expected 2 active reactions after two users add, got %d", activeTotal)
	}

	removeUserOneRes := s.HandleReaction(ctx, userOne, commentID, &commentrequest.AddReactionRequest{Type: "like"})
	if !removeUserOneRes.Success {
		t.Fatalf("expected user one remove reaction success, got: %s", removeUserOneRes.Message)
	}

	if err := db.Raw("SELECT COUNT(*) FROM comment_reactions WHERE comment_id = ? AND deleted_at IS NULL", commentID).Scan(&activeTotal).Error; err != nil {
		t.Fatalf("failed to count active reactions after user one remove: %v", err)
	}
	if activeTotal != 1 {
		t.Fatalf("expected 1 active reaction after user one remove, got %d", activeTotal)
	}

	var userTwoActive int64
	if err := db.Raw("SELECT COUNT(*) FROM comment_reactions WHERE comment_id = ? AND user_id = ? AND deleted_at IS NULL", commentID, userTwoID).Scan(&userTwoActive).Error; err != nil {
		t.Fatalf("failed to count user two active reactions: %v", err)
	}
	if userTwoActive != 1 {
		t.Fatalf("expected user two reaction to remain active, got %d", userTwoActive)
	}
}
