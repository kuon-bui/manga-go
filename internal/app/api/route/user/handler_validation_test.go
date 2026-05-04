package userroute

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func newUserCtx(method, path string, body string) (*gin.Context, *httptest.ResponseRecorder) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req := httptest.NewRequest(method, path, bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	c.Request = req
	return c, w
}

func TestCreateAccountInvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := &userHandler{}
	c, w := newUserCtx(http.MethodPost, "/users", "{")

	h.createAccount(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}

func TestSignInInvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := &userHandler{}
	c, w := newUserCtx(http.MethodPost, "/users/sign-in", "{")

	h.signIn(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}

}

func TestRequestResetPasswordInvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := &userHandler{}
	c, w := newUserCtx(http.MethodPost, "/users/request-reset-password", "{")

	h.requestResetPassword(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}

func TestResetPasswordInvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := &userHandler{}
	c, w := newUserCtx(http.MethodPost, "/users/reset-password", "{")

	h.resetPassword(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}

func TestGetUserRolesInvalidID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := &userHandler{}
	c, w := newUserCtx(http.MethodGet, "/users/invalid/roles", "")
	c.Params = gin.Params{{Key: "id", Value: "invalid"}}

	h.getUserRoles(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}

func TestAssignUserRoleInvalidID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := &userHandler{}
	c, w := newUserCtx(http.MethodPost, "/users/invalid/roles", "{}")
	c.Params = gin.Params{{Key: "id", Value: "invalid"}}

	h.assignUserRole(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}

func TestAssignUserRoleInvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := &userHandler{}
	c, w := newUserCtx(http.MethodPost, "/users/550e8400-e29b-41d4-a716-446655440000/roles", "{")
	c.Params = gin.Params{{Key: "id", Value: "550e8400-e29b-41d4-a716-446655440000"}}

	h.assignUserRole(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}

func TestRemoveUserRoleInvalidUserID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := &userHandler{}
	c, w := newUserCtx(http.MethodDelete, "/users/invalid/roles/role", "")
	c.Params = gin.Params{{Key: "id", Value: "invalid"}, {Key: "roleId", Value: "550e8400-e29b-41d4-a716-446655440000"}}

	h.removeUserRole(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}

func TestRemoveUserRoleInvalidRoleID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := &userHandler{}
	c, w := newUserCtx(http.MethodDelete, "/users/user/roles/invalid", "")
	c.Params = gin.Params{{Key: "id", Value: "550e8400-e29b-41d4-a716-446655440000"}, {Key: "roleId", Value: "invalid"}}

	h.removeUserRole(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}
