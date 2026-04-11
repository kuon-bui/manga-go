package authmiddleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func newGinContext(method string, cookies ...*http.Cookie) (*gin.Context, *httptest.ResponseRecorder) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	req := httptest.NewRequest(method, "/", nil)
	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}

	c.Request = req
	return c, w
}

func TestExtractTokenFromCookiesSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)

	c, w := newGinContext(http.MethodGet, &http.Cookie{Name: "access_token", Value: "token-1"})
	m := &AuthMiddleware{}

	token := m.extractTokenFromCookies(c, "access_token")

	if token != "token-1" {
		t.Fatalf("expected token-1, got %s", token)
	}
	if c.IsAborted() {
		t.Fatalf("context should not be aborted for valid cookie")
	}
	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}
}

func TestExtractTokenFromCookiesMissingCookie(t *testing.T) {
	gin.SetMode(gin.TestMode)

	c, w := newGinContext(http.MethodGet)
	m := &AuthMiddleware{}

	token := m.extractTokenFromCookies(c, "access_token")

	if token != "" {
		t.Fatalf("expected empty token, got %s", token)
	}
	if !c.IsAborted() {
		t.Fatalf("context should be aborted when cookie is missing")
	}
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", w.Code)
	}
}

func TestExtractTokenFromCookiesEmptyValue(t *testing.T) {
	gin.SetMode(gin.TestMode)

	c, w := newGinContext(http.MethodGet, &http.Cookie{Name: "access_token", Value: ""})
	m := &AuthMiddleware{}

	token := m.extractTokenFromCookies(c, "access_token")

	if token != "" {
		t.Fatalf("expected empty token, got %s", token)
	}
	if !c.IsAborted() {
		t.Fatalf("context should be aborted when token is empty")
	}
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", w.Code)
	}
}
