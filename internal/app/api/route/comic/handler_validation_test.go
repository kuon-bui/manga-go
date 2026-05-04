package comicroute

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func newComicCtx(method, path string, body string) (*gin.Context, *httptest.ResponseRecorder) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req := httptest.NewRequest(method, path, bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	c.Request = req
	return c, w
}

func TestCreateComicInvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := &ComicHandler{}
	c, w := newComicCtx(http.MethodPost, "/comics", "{")

	h.createComic(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}

func TestGetComicsInvalidQuery(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := &ComicHandler{}
	c, w := newComicCtx(http.MethodGet, "/comics?page=invalid", "")

	h.getComics(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}

func TestUpdateComicInvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := &ComicHandler{}
	c, w := newComicCtx(http.MethodPut, "/comics/my-comic", "{")
	c.Params = gin.Params{{Key: "comicSlug", Value: "my-comic"}}

	h.updateComic(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}

func TestUpdateComicStatusInvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := &ComicHandler{}
	c, w := newComicCtx(http.MethodPatch, "/comics/my-comic/status", "{")
	c.Params = gin.Params{{Key: "comicSlug", Value: "my-comic"}}

	h.updateComicStatus(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}

func TestPublishComicInvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := &ComicHandler{}
	c, w := newComicCtx(http.MethodPatch, "/comics/my-comic/publish", "{")
	c.Params = gin.Params{{Key: "comicSlug", Value: "my-comic"}}

	h.publishComic(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}
