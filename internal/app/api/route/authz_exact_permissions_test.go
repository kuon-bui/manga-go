package route_test

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	authorizationroute "manga-go/internal/app/api/route/authorization"
	permissionroute "manga-go/internal/app/api/route/permission"
	roleroute "manga-go/internal/app/api/route/role"
	userroute "manga-go/internal/app/api/route/user"
	authmiddleware "manga-go/internal/app/middleware/auth"
	authzmiddleware "manga-go/internal/app/middleware/authz"
	"manga-go/internal/pkg/authorization"
	"manga-go/internal/pkg/common"
	"manga-go/internal/pkg/config"
	jwtprovider "manga-go/internal/pkg/jwt_provider"
	"manga-go/internal/pkg/logger"
	"manga-go/internal/pkg/model"
	redisstore "manga-go/internal/pkg/redis"
	userrepo "manga-go/internal/pkg/repo/user"
	"manga-go/internal/pkg/testutil"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	goredis "github.com/redis/go-redis/v9"
)

type exactRouteCase struct {
	name       string
	method     string
	path       string
	permission string
	register   func(*gin.Engine, *authmiddleware.AuthMiddleware, *authzmiddleware.AuthzMiddleware)
}

func TestAuthorizationAdminRoutesRequireExactPermission(t *testing.T) {
	auth, authz, policy, roleID, token := newExactPermissionTestAuth(t)
	cases := []exactRouteCase{
		{"users list", http.MethodGet, "/users", "user:read", registerExactUserRoutes},
		{"user authorization", http.MethodGet, "/users/11111111-1111-1111-1111-111111111111/authorization", "user:read", registerExactUserRoutes},
		{"replace user roles", http.MethodPost, "/users/11111111-1111-1111-1111-111111111111/roles", "role:manage", registerExactUserRoutes},
		{"remove user role", http.MethodDelete, "/users/11111111-1111-1111-1111-111111111111/roles/22222222-2222-2222-2222-222222222222", "role:manage", registerExactUserRoutes},
		{"roles list", http.MethodGet, "/roles/all", "role:manage", registerExactRoleRoutes},
		{"role detail", http.MethodGet, "/roles/11111111-1111-1111-1111-111111111111", "role:manage", registerExactRoleRoutes},
		{"create role", http.MethodPost, "/roles", "role:manage", registerExactRoleRoutes},
		{"update role", http.MethodPut, "/roles/11111111-1111-1111-1111-111111111111", "role:manage", registerExactRoleRoutes},
		{"delete role", http.MethodDelete, "/roles/11111111-1111-1111-1111-111111111111", "role:manage", registerExactRoleRoutes},
		{"replace role permissions", http.MethodPost, "/roles/11111111-1111-1111-1111-111111111111/permissions", "role:manage", registerExactRoleRoutes},
		{"remove role permission", http.MethodDelete, "/roles/11111111-1111-1111-1111-111111111111/permissions/role:manage", "role:manage", registerExactRoleRoutes},
		{"permission catalog", http.MethodGet, "/permissions", "permission:read", registerExactPermissionRoutes},
		{"authorization audit", http.MethodGet, "/authorization/audit-logs", "audit_log:read", registerExactAuthorizationRoutes},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if err := policy.ReplacePermissionsForRole(
				roleID.String(), []string{"comic:read"}, authorization.OrgPlatform,
			); err != nil {
				t.Fatal(err)
			}
			if status := serveExactPermissionRoute(tc, auth, authz, token); status != http.StatusForbidden {
				t.Fatalf("wrong permission reached %s %s: status %d", tc.method, tc.path, status)
			}

			if err := policy.ReplacePermissionsForRole(
				roleID.String(), []string{tc.permission}, authorization.OrgPlatform,
			); err != nil {
				t.Fatal(err)
			}
			if status := serveExactPermissionRoute(tc, auth, authz, token); status == http.StatusForbidden || status == http.StatusUnauthorized {
				t.Fatalf("exact permission %s did not pass %s %s: status %d", tc.permission, tc.method, tc.path, status)
			}
		})
	}
}

func serveExactPermissionRoute(
	tc exactRouteCase,
	auth *authmiddleware.AuthMiddleware,
	authz *authzmiddleware.AuthzMiddleware,
	token string,
) (status int) {
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	tc.register(engine, auth, authz)
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(tc.method, tc.path, strings.NewReader("{}"))
	request.Header.Set("Content-Type", "application/json")
	request.AddCookie(&http.Cookie{Name: "access_token", Value: token})

	defer func() {
		if recover() != nil {
			status = recorder.Code
		}
	}()
	engine.ServeHTTP(recorder, request)
	return recorder.Code
}

func newExactPermissionTestAuth(
	t *testing.T,
) (*authmiddleware.AuthMiddleware, *authzmiddleware.AuthzMiddleware, *authorization.PolicyManager, uuid.UUID, string) {
	t.Helper()
	db := testutil.NewSQLiteDB(t)
	testutil.MustSyncSchemas(t, db, &testutil.User{})
	user := &model.User{
		SqlModel: common.SqlModel{ID: uuid.New()},
		Name:     "route-admin",
		Email:    "route-admin@example.com",
	}
	if err := db.Create(user).Error; err != nil {
		t.Fatal(err)
	}

	configPath := filepath.Join(t.TempDir(), "config.yml")
	configYAML := `
run_mode: development
jwt:
  secret: route-test-secret
  refresh_secret: route-test-refresh-secret
  expires_seconds: 3600
  refresh_expire_seconds: 7200
cookie_name:
  access_token: access_token
  refresh_token: refresh_token
`
	if err := os.WriteFile(configPath, []byte(configYAML), 0o600); err != nil {
		t.Fatal(err)
	}
	cfg := config.LoadConfig(configPath)()
	redisClient := newNilRedisClient(t)
	redisWrapper := redisstore.NewRedis(redisClient, logger.NewLogger())
	jwt := jwtprovider.NewJwtProvider(cfg, logger.NewLogger(), redisWrapper)
	token, _, err := jwt.GenerateToken(jwtprovider.UserPayload{
		UserID: user.ID, FullName: user.Name, Email: user.Email,
	})
	if err != nil {
		t.Fatal(err)
	}
	auth := authmiddleware.NewAuthMiddleware(authmiddleware.AuthMiddlewareParams{
		Jwt: jwt, Config: cfg, UserRepo: userrepo.NewUserRepository(db, nil),
	})

	enforcer := testutil.NewInMemoryEnforcer(t)
	policy := authorization.NewPolicyManager(authorization.PolicyManagerParams{Enforcer: enforcer})
	roleID := uuid.New()
	if err := policy.AddRoleForUser(user.ID.String(), roleID.String(), authorization.OrgPlatform); err != nil {
		t.Fatal(err)
	}
	authz := authzmiddleware.NewAuthzMiddleware(authzmiddleware.AuthzMiddlewareParams{
		Authorizer: authorization.NewAuthorizer(enforcer),
	})
	return auth, authz, policy, roleID, token.TokenString
}

func newNilRedisClient(t *testing.T) *goredis.Client {
	t.Helper()
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = listener.Close() })
	go func() {
		for {
			conn, acceptErr := listener.Accept()
			if acceptErr != nil {
				return
			}
			go serveNilRedisConnection(conn)
		}
	}()
	client := goredis.NewClient(&goredis.Options{
		Addr: listener.Addr().String(), Protocol: 2, DisableIdentity: true,
	})
	t.Cleanup(func() { _ = client.Close() })
	return client
}

func serveNilRedisConnection(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return
		}
		count, err := strconv.Atoi(strings.TrimSpace(strings.TrimPrefix(line, "*")))
		if err != nil {
			return
		}
		for range count {
			lengthLine, err := reader.ReadString('\n')
			if err != nil {
				return
			}
			length, err := strconv.Atoi(strings.TrimSpace(strings.TrimPrefix(lengthLine, "$")))
			if err != nil {
				return
			}
			if _, err := io.CopyN(io.Discard, reader, int64(length+2)); err != nil {
				return
			}
		}
		if _, err := fmt.Fprint(conn, "$-1\r\n"); err != nil {
			return
		}
	}
}

func registerExactUserRoutes(e *gin.Engine, auth *authmiddleware.AuthMiddleware, authz *authzmiddleware.AuthzMiddleware) {
	userroute.NewUserRoute(userroute.UserRouteParams{
		R: e, UserHandler: userroute.NewUserHandler(userroute.UserHandlerParams{}),
		AuthMiddleware: auth, AuthzMiddleware: authz,
	}).Setup()
}

func registerExactRoleRoutes(e *gin.Engine, auth *authmiddleware.AuthMiddleware, authz *authzmiddleware.AuthzMiddleware) {
	roleroute.NewRoleRoute(roleroute.RoleRouteParams{
		R: e, RoleHandler: roleroute.NewRoleHandler(roleroute.RoleHandlerParams{}),
		AuthMiddleware: auth, AuthzMiddleware: authz,
	}).Setup()
}

func registerExactPermissionRoutes(e *gin.Engine, auth *authmiddleware.AuthMiddleware, authz *authzmiddleware.AuthzMiddleware) {
	permissionroute.NewPermissionRoute(permissionroute.PermissionRouteParams{
		R: e, PermissionHandler: permissionroute.NewPermissionHandler(permissionroute.PermissionHandlerParams{}),
		AuthMiddleware: auth, AuthzMiddleware: authz,
	}).Setup()
}

func registerExactAuthorizationRoutes(e *gin.Engine, auth *authmiddleware.AuthMiddleware, authz *authzmiddleware.AuthzMiddleware) {
	authorizationroute.NewRoute(authorizationroute.RouteParams{
		Engine: e, Handler: authorizationroute.NewHandler(authorizationroute.HandlerParams{}),
		AuthMiddleware: auth, AuthzMiddleware: authz,
	}).Setup()
}
