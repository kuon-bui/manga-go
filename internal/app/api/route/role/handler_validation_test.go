package roleroute

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func newRoleCtx(method, path string, body string) (*gin.Context, *httptest.ResponseRecorder) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req := httptest.NewRequest(method, path, bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	c.Request = req
	return c, w
}

func TestCreateRoleInvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := &RoleHandler{}
	c, w := newRoleCtx(http.MethodPost, "/roles", "{")

	h.createRole(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}

func TestGetRolesInvalidQuery(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := &RoleHandler{}
	c, w := newRoleCtx(http.MethodGet, "/roles?page=invalid", "")

	h.getRoles(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}

func TestGetRoleInvalidID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := &RoleHandler{}
	c, w := newRoleCtx(http.MethodGet, "/roles/invalid", "")
	c.Params = gin.Params{{Key: "id", Value: "invalid"}}

	h.getRole(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}

func TestUpdateRoleInvalidID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := &RoleHandler{}
	c, w := newRoleCtx(http.MethodPut, "/roles/invalid", "{}")
	c.Params = gin.Params{{Key: "id", Value: "invalid"}}

	h.updateRole(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}

func TestUpdateRoleInvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := &RoleHandler{}
	c, w := newRoleCtx(http.MethodPut, "/roles/550e8400-e29b-41d4-a716-446655440000", "{")
	c.Params = gin.Params{{Key: "id", Value: "550e8400-e29b-41d4-a716-446655440000"}}

	h.updateRole(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}

func TestDeleteRoleInvalidID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := &RoleHandler{}
	c, w := newRoleCtx(http.MethodDelete, "/roles/invalid", "")
	c.Params = gin.Params{{Key: "id", Value: "invalid"}}

	h.deleteRole(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}

func TestAssignRolePermissionInvalidID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := &RoleHandler{}
	c, w := newRoleCtx(http.MethodPost, "/roles/invalid/permissions", "{}")
	c.Params = gin.Params{{Key: "id", Value: "invalid"}}

	h.assignRolePermission(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}

func TestAssignRolePermissionInvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := &RoleHandler{}
	c, w := newRoleCtx(http.MethodPost, "/roles/550e8400-e29b-41d4-a716-446655440000/permissions", "{")
	c.Params = gin.Params{{Key: "id", Value: "550e8400-e29b-41d4-a716-446655440000"}}

	h.assignRolePermission(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}

func TestRemoveRolePermissionInvalidRoleID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := &RoleHandler{}
	c, w := newRoleCtx(http.MethodDelete, "/roles/invalid/permissions/p", "")
	c.Params = gin.Params{{Key: "id", Value: "invalid"}, {Key: "permissionId", Value: "550e8400-e29b-41d4-a716-446655440000"}}

	h.removeRolePermission(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}

func TestRemoveRolePermissionInvalidPermissionID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := &RoleHandler{}
	c, w := newRoleCtx(http.MethodDelete, "/roles/id/permissions/invalid", "")
	c.Params = gin.Params{{Key: "id", Value: "550e8400-e29b-41d4-a716-446655440000"}, {Key: "permissionId", Value: "invalid"}}

	h.removeRolePermission(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}
