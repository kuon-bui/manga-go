package commentroute

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"manga-go/internal/pkg/common"
	"manga-go/internal/pkg/model"

	"github.com/gin-gonic/gin"
)

func newCommentCtx(method, path string, body string) (*gin.Context, *httptest.ResponseRecorder) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req := httptest.NewRequest(method, path, bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	c.Request = req
	return c, w
}

func TestCreateCommentUnauthorizedWithoutCurrentUser(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := &CommentHandler{}
	c, w := newCommentCtx(http.MethodPost, "/comments", "")
	c.Set(string(common.CurrentUser), (*model.User)(nil))

	h.createComment(c)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", w.Code)
	}
}

func TestGetCommentInvalidID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := &CommentHandler{}
	c, w := newCommentCtx(http.MethodGet, "/comments/invalid", "")
	c.Params = gin.Params{{Key: "id", Value: "invalid"}}

	h.getComment(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}

func TestDeleteCommentInvalidID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := &CommentHandler{}
	c, w := newCommentCtx(http.MethodDelete, "/comments/invalid", "")
	c.Params = gin.Params{{Key: "id", Value: "invalid"}}

	h.deleteComment(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}

func TestHandleReactionInvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := &CommentHandler{}
	c, w := newCommentCtx(http.MethodPost, "/comments/x/reactions", "{")
	c.Params = gin.Params{{Key: "id", Value: "550e8400-e29b-41d4-a716-446655440000"}}

	h.handleReaction(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}

func TestHandleReactionInvalidID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := &CommentHandler{}
	c, w := newCommentCtx(http.MethodPost, "/comments/invalid/reactions", `{"type":"like"}`)
	c.Params = gin.Params{{Key: "id", Value: "invalid"}}

	h.handleReaction(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}
