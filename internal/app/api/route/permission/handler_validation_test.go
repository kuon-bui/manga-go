package permissionroute

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func newPermissionCtx(method, path string, body string) (*gin.Context, *httptest.ResponseRecorder) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req := httptest.NewRequest(method, path, bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	c.Request = req
	return c, w
}

func TestCreatePermissionInvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := &PermissionHandler{}
	c, w := newPermissionCtx(http.MethodPost, "/permissions", "{")

	h.createPermission(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}

func TestGetPermissionsInvalidQuery(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := &PermissionHandler{}
	c, w := newPermissionCtx(http.MethodGet, "/permissions?page=invalid", "")

	h.getPermissions(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}

func TestUpdatePermissionInvalidID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := &PermissionHandler{}
	c, w := newPermissionCtx(http.MethodPut, "/permissions/invalid", "{}")
	c.Params = gin.Params{{Key: "id", Value: "invalid"}}

	h.updatePermission(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}

func TestUpdatePermissionInvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := &PermissionHandler{}
	c, w := newPermissionCtx(http.MethodPut, "/permissions/550e8400-e29b-41d4-a716-446655440000", "{")
	c.Params = gin.Params{{Key: "id", Value: "550e8400-e29b-41d4-a716-446655440000"}}

	h.updatePermission(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}

func TestDeletePermissionInvalidID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := &PermissionHandler{}
	c, w := newPermissionCtx(http.MethodDelete, "/permissions/invalid", "")
	c.Params = gin.Params{{Key: "id", Value: "invalid"}}

	h.deletePermission(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}
