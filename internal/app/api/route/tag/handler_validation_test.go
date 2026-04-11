package tagroute

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func newTagCtx(method, path string, body string) (*gin.Context, *httptest.ResponseRecorder) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req := httptest.NewRequest(method, path, bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	c.Request = req
	return c, w
}

func TestCreateTagInvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := &TagHandler{}
	c, w := newTagCtx(http.MethodPost, "/tags", "{")

	h.createTag(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}

func TestGetTagsInvalidQuery(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := &TagHandler{}
	c, w := newTagCtx(http.MethodGet, "/tags?page=invalid", "")

	h.getTags(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}
