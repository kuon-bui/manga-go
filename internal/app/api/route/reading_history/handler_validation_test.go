package readinghistoryroute

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"manga-go/internal/pkg/common"
	"manga-go/internal/pkg/model"

	"github.com/gin-gonic/gin"
)

func newReadingHistoryCtx(method, path string) (*gin.Context, *httptest.ResponseRecorder) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(method, path, nil)
	return c, w
}

func newReadingHistoryCtxWithBody(method, path string, body string) (*gin.Context, *httptest.ResponseRecorder) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req := httptest.NewRequest(method, path, bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	c.Request = req
	return c, w
}

func TestCreateReadingHistoryUnauthorizedWithoutCurrentUser(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := &ReadingHistoryHandler{}
	c, w := newReadingHistoryCtx(http.MethodPost, "/reading-histories")
	c.Set(string(common.CurrentUser), (*model.User)(nil))

	h.createReadingHistory(c)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", w.Code)
	}
}

func TestGetReadingHistoryInvalidID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := &ReadingHistoryHandler{}
	c, w := newReadingHistoryCtx(http.MethodGet, "/reading-histories/invalid")
	c.Params = gin.Params{{Key: "id", Value: "invalid"}}

	h.getReadingHistory(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}

func TestGetReadingHistoriesInvalidQuery(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := &ReadingHistoryHandler{}
	c, w := newReadingHistoryCtx(http.MethodGet, "/reading-histories?page=invalid")
	c.Set(string(common.CurrentUser), &model.User{})

	h.getReadingHistories(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}

func TestCreateReadingHistoryInvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := &ReadingHistoryHandler{}
	c, w := newReadingHistoryCtxWithBody(http.MethodPost, "/reading-histories", "{")
	c.Set(string(common.CurrentUser), &model.User{})

	h.createReadingHistory(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}

func TestUpdateReadingHistoryInvalidID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := &ReadingHistoryHandler{}
	c, w := newReadingHistoryCtx(http.MethodPut, "/reading-histories/invalid")
	c.Params = gin.Params{{Key: "id", Value: "invalid"}}

	h.updateReadingHistory(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}

func TestUpdateReadingHistoryInvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := &ReadingHistoryHandler{}
	c, w := newReadingHistoryCtxWithBody(http.MethodPut, "/reading-histories/550e8400-e29b-41d4-a716-446655440000", "{")
	c.Params = gin.Params{{Key: "id", Value: "550e8400-e29b-41d4-a716-446655440000"}}

	h.updateReadingHistory(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}

func TestDeleteReadingHistoryInvalidID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := &ReadingHistoryHandler{}
	c, w := newReadingHistoryCtx(http.MethodDelete, "/reading-histories/invalid")
	c.Params = gin.Params{{Key: "id", Value: "invalid"}}

	h.deleteReadingHistory(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}
