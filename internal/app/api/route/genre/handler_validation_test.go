package genreroute

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func newGenreCtx(method, path string, body string) (*gin.Context, *httptest.ResponseRecorder) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req := httptest.NewRequest(method, path, bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	c.Request = req
	return c, w
}

func TestCreateGenreInvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := &GenreHandler{}
	c, w := newGenreCtx(http.MethodPost, "/genres", "{")

	h.createGenre(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}

func TestGetGenresInvalidQuery(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := &GenreHandler{}
	c, w := newGenreCtx(http.MethodGet, "/genres?page=invalid", "")

	h.getGenres(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}

func TestUpdateGenreInvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := &GenreHandler{}
	c, w := newGenreCtx(http.MethodPut, "/genres/action", "{")
	c.Params = gin.Params{{Key: "slug", Value: "action"}}

	h.updateGenre(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}
