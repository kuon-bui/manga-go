package authorroute

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func newAuthorCtx(method, path string, body string) (*gin.Context, *httptest.ResponseRecorder) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req := httptest.NewRequest(method, path, bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	c.Request = req
	return c, w
}

func TestCreateAuthorInvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := &AuthorHandler{}
	c, w := newAuthorCtx(http.MethodPost, "/authors", "{")

	h.createAuthor(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}

func TestGetAuthorsInvalidQuery(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := &AuthorHandler{}
	c, w := newAuthorCtx(http.MethodGet, "/authors?page=invalid", "")

	h.getAuthors(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}

func TestUpdateAuthorInvalidID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := &AuthorHandler{}
	c, w := newAuthorCtx(http.MethodPut, "/authors/invalid", "{}")
	c.Params = gin.Params{{Key: "id", Value: "invalid"}}

	h.updateAuthor(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}

func TestUpdateAuthorInvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := &AuthorHandler{}
	c, w := newAuthorCtx(http.MethodPut, "/authors/550e8400-e29b-41d4-a716-446655440000", "{")
	c.Params = gin.Params{{Key: "id", Value: "550e8400-e29b-41d4-a716-446655440000"}}

	h.updateAuthor(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}
