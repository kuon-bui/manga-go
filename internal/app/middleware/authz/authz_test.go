package authzmiddleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"manga-go/internal/pkg/authorization"

	"github.com/gin-gonic/gin"
)

func newGinContext() (*gin.Context, *httptest.ResponseRecorder) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	return c, w
}

func TestRequireDeniesWhenMiddlewareIsNil(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := Require(nil, authorization.ActionManage, authorization.ObjectRole)
	c, w := newGinContext()

	handler(c)

	if !c.IsAborted() {
		t.Fatal("expected the request to be aborted when the middleware is missing, it was allowed through")
	}
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500 for a missing middleware, got %d", w.Code)
	}
}

func TestRequireDeniesWhenAuthorizerIsMissing(t *testing.T) {
	gin.SetMode(gin.TestMode)

	m := &AuthzMiddleware{}
	handler := m.Require(authorization.ActionManage, authorization.ObjectRole)
	c, w := newGinContext()

	handler(c)

	if !c.IsAborted() {
		t.Fatal("expected the request to be aborted when the authorizer is missing, it was allowed through")
	}
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500 for a missing authorizer, got %d", w.Code)
	}
}

func TestRequireDeniesWhenNoUserIsInContext(t *testing.T) {
	gin.SetMode(gin.TestMode)

	m := &AuthzMiddleware{authorizer: authorization.NewAuthorizer(nil)}
	handler := m.Require(authorization.ActionManage, authorization.ObjectRole)
	c, w := newGinContext()

	handler(c)

	if !c.IsAborted() {
		t.Fatal("expected the request to be aborted when no user is in context, it was allowed through")
	}
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 when no user is in context, got %d", w.Code)
	}
}
