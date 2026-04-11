package slugmiddleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func newTestContext() (*gin.Context, *httptest.ResponseRecorder) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	return c, w
}

func TestNewSlugMiddleware(t *testing.T) {
	if m := NewSlugMiddleware(SlugMiddlewareParams{}); m == nil {
		t.Fatal("expected slug middleware instance")
	}
}

func TestResolveComicIDRequiresSlug(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, w := newTestContext()

	(&SlugMiddleware{}).ResolveComicID(c)

	if !c.IsAborted() {
		t.Fatal("expected context to be aborted")
	}
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}

func TestResolveChapterIDRequiresSlug(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, w := newTestContext()

	(&SlugMiddleware{}).ResolveChapterID(c)

	if !c.IsAborted() {
		t.Fatal("expected context to be aborted")
	}
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}

func TestResolveTranslationGroupIDRequiresSlug(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, w := newTestContext()

	(&SlugMiddleware{}).ResolveTranslationGroupID(c)

	if !c.IsAborted() {
		t.Fatal("expected context to be aborted")
	}
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}

func TestResolveGenreIDRequiresSlug(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, w := newTestContext()

	(&SlugMiddleware{}).ResolveGenreID(c)

	if !c.IsAborted() {
		t.Fatal("expected context to be aborted")
	}
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}

func TestResolveTagIDRequiresSlug(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, w := newTestContext()

	(&SlugMiddleware{}).ResolveTagID(c)

	if !c.IsAborted() {
		t.Fatal("expected context to be aborted")
	}
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}
