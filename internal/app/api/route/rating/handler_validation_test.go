package ratingroute

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"manga-go/internal/pkg/common"
	"manga-go/internal/pkg/model"

	"github.com/gin-gonic/gin"
)

func newRatingCtx(method, path string) (*gin.Context, *httptest.ResponseRecorder) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(method, path, nil)
	return c, w
}

func TestCreateRatingUnauthorizedWithoutCurrentUser(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := &RatingHandler{}
	c, w := newRatingCtx(http.MethodPost, "/ratings/comics/x")
	c.Set(string(common.CurrentUser), (*model.User)(nil))

	h.createRating(c)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", w.Code)
	}
}

func TestUpdateRatingInvalidID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := &RatingHandler{}
	c, w := newRatingCtx(http.MethodPut, "/ratings/comics/x/invalid")
	c.Set(string(common.CurrentUser), &model.User{})
	c.Params = gin.Params{{Key: "id", Value: "invalid"}}

	h.updateRating(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}

func TestDeleteRatingInvalidID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := &RatingHandler{}
	c, w := newRatingCtx(http.MethodDelete, "/ratings/comics/x/invalid")
	c.Set(string(common.CurrentUser), &model.User{})
	c.Params = gin.Params{{Key: "id", Value: "invalid"}}

	h.deleteRating(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}
