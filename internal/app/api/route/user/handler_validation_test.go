package userroute

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"manga-go/internal/pkg/model"
	authorizationrevision "manga-go/internal/pkg/repo/authorization_revision"
	userrepo "manga-go/internal/pkg/repo/user"
	userrequest "manga-go/internal/pkg/request/user"
	authorizationadmin "manga-go/internal/pkg/services/authorization_admin"
	"manga-go/internal/pkg/testutil"

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

func TestGetMyAuthorizationRequiresAuthenticatedViewer(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := &userHandler{}
	c, w := newUserCtx(http.MethodGet, "/users/me/authorization", "")

	h.getMyAuthorization(c)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", w.Code)
	}
}

func TestGetUsersRejectsInvalidRoleFilter(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := &userHandler{}
	c, w := newUserCtx(http.MethodGet, "/users?role_id=not-a-uuid", "")
	c.Request.URL.RawQuery = "role_id=not-a-uuid"

	h.getUsers(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}

func TestGetUsersRejectsLimitAboveMaximum(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := &userHandler{}
	c, w := newUserCtx(http.MethodGet, "/users?limit=101", "")
	c.Request.URL.RawQuery = "limit=101"

	h.getUsers(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}

func TestGetUsersUsesStandardPaginationResponse(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := testutil.NewSQLiteDB(t)
	testutil.MustSyncSchemas(t, db, &testutil.User{}, &model.AuthorizationCacheRevision{})
	h := &userHandler{authorizationAdmin: authorizationadmin.NewService(authorizationadmin.ServiceParams{
		UserRepo:  userrepo.NewUserRepository(db, nil),
		Revisions: authorizationrevision.NewRepo(db),
	})}
	c, w := newUserCtx(http.MethodGet, "/users", "")

	h.getUsers(c)

	var body struct {
		Message string                     `json:"message"`
		Data    map[string]json.RawMessage `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if body.Message != "Users retrieved successfully" {
		t.Fatalf("unexpected message: %q", body.Message)
	}
	if len(body.Data) != 2 {
		t.Fatalf("expected only data and total pagination fields, got %s", body.Data)
	}
	var users []json.RawMessage
	if err := json.Unmarshal(body.Data["data"], &users); err != nil {
		t.Fatalf("expected user data array: %v", err)
	}
	var total int64
	if err := json.Unmarshal(body.Data["total"], &total); err != nil {
		t.Fatalf("expected user total: %v", err)
	}
	if len(users) != 0 || total != 0 {
		t.Fatalf("unexpected user page: users=%d total=%d", len(users), total)
	}
}

func TestGetAuthorizationUserRejectsInvalidID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := &userHandler{}
	c, w := newUserCtx(http.MethodGet, "/users/invalid/authorization", "")
	c.Params = gin.Params{{Key: "id", Value: "invalid"}}

	h.getAuthorizationUser(c)

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

func TestAssignUserRoleRequestAcceptsEmptyRoleSet(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := newUserCtx(http.MethodPost, "/users/id/roles", `{"role_ids":[]}`)
	var request userrequest.AssignRoleRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		t.Fatalf("expected an empty replacement set to bind: %v", err)
	}
	if request.RoleIDs == nil || len(*request.RoleIDs) != 0 {
		t.Fatalf("expected a present empty role set, got %#v", request.RoleIDs)
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
