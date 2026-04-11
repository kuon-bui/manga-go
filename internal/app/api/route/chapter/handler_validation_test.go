package chapterhandler

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"manga-go/internal/pkg/common"
	"manga-go/internal/pkg/model"

	"github.com/gin-gonic/gin"
)

func newChapterCtx(method, path string, body string) (*gin.Context, *httptest.ResponseRecorder) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req := httptest.NewRequest(method, path, bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	c.Request = req
	return c, w
}

func TestCreateChapterInvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := &ChapterHandler{}
	c, w := newChapterCtx(http.MethodPost, "/comics/s/chapters", "{")

	h.createChapter(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}

func TestListChaptersInvalidQuery(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := &ChapterHandler{}
	c, w := newChapterCtx(http.MethodGet, "/comics/s/chapters?page=invalid", "")

	h.listChapters(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}

func TestUpdateChapterInvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := &ChapterHandler{}
	c, w := newChapterCtx(http.MethodPut, "/comics/s/chapters/ch", "{")
	c.Params = gin.Params{{Key: "chapterSlug", Value: "ch"}}

	h.updateChapter(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}

func TestPublishChapterInvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := &ChapterHandler{}
	c, w := newChapterCtx(http.MethodPatch, "/comics/s/chapters/ch/publish", "{")
	c.Params = gin.Params{{Key: "chapterSlug", Value: "ch"}}

	h.publishChapter(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}

func TestUpdateChapterPagesInvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := &ChapterHandler{}
	c, w := newChapterCtx(http.MethodPut, "/comics/s/chapters/ch/pages", "{")
	c.Params = gin.Params{{Key: "chapterSlug", Value: "ch"}}

	h.updateChapterPages(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}

func TestMarkChapterAsReadWithoutContextIDs(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := &ChapterHandler{}
	c, w := newChapterCtx(http.MethodPatch, "/comics/s/chapters/ch/mark-as-read", "")

	h.markChapterAsRead(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}

func TestGetReadingProgressUnauthorizedWithoutCurrentUser(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := &ChapterHandler{}
	c, w := newChapterCtx(http.MethodGet, "/comics/s/chapters/ch/reading-progress", "")
	c.Set(string(common.CurrentUser), (*model.User)(nil))

	h.getReadingProgress(c)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", w.Code)
	}
}

func TestUpdateReadingProgressInvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := &ChapterHandler{}
	c, w := newChapterCtx(http.MethodPatch, "/comics/s/chapters/ch/reading-progress", "{")
	c.Set(string(common.CurrentUser), &model.User{})

	h.updateReadingProgress(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}
