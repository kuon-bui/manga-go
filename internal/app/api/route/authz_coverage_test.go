package route_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	authorizationroute "manga-go/internal/app/api/route/authorization"
	comicstatsroute "manga-go/internal/app/api/route/comic_stats"
	permissionroute "manga-go/internal/app/api/route/permission"
	roleroute "manga-go/internal/app/api/route/role"
	userroute "manga-go/internal/app/api/route/user"
	authmiddleware "manga-go/internal/app/middleware/auth"
	authzmiddleware "manga-go/internal/app/middleware/authz"

	"github.com/gin-gonic/gin"
)

// handlerChainFor returns the names of every handler gin would run for a
// request, including middleware. gin only exposes the final handler through
// Routes(), but a middleware registered on the engine before the routes runs
// first and can read the whole matched chain.
func handlerChainFor(t *testing.T, register func(e *gin.Engine), method string, path string) []string {
	t.Helper()
	gin.SetMode(gin.TestMode)

	e := gin.New()
	var chain []string
	e.Use(func(c *gin.Context) {
		chain = c.HandlerNames()
		c.AbortWithStatus(http.StatusTeapot)
	})
	register(e)

	w := httptest.NewRecorder()
	e.ServeHTTP(w, httptest.NewRequest(method, path, nil))

	if w.Code != http.StatusTeapot {
		t.Fatalf("expected the probe middleware to run for %s %s, got status %d", method, path, w.Code)
	}
	return chain
}

// The guard closure gets inlined into the caller's package, so the symbol is
// named after the route package rather than middleware/authz. What stays stable
// is the method it came from.
var authzGuardMarkers = []string{"AuthzMiddleware).Require", "authz.denyUnavailable"}

func assertGuardedByAuthz(t *testing.T, chain []string, method string, path string) {
	t.Helper()

	for _, name := range chain {
		for _, marker := range authzGuardMarkers {
			if strings.Contains(name, marker) {
				return
			}
		}
	}
	t.Fatalf("%s %s runs without an authorization guard; chain was %v", method, path, chain)
}

func registerComicStatsRoutes(e *gin.Engine) {
	comicstatsroute.NewComicStatsRoute(comicstatsroute.ComicStatsRouteParams{
		R:               e,
		AuthMiddleware:  &authmiddleware.AuthMiddleware{},
		AuthzMiddleware: &authzmiddleware.AuthzMiddleware{},
		Handler:         comicstatsroute.NewComicStatsHandler(comicstatsroute.ComicStatsHandlerParams{}),
	}).Setup()
}

// These endpoints recompute stats across the catalogue. Authentication alone is
// not enough: any signed-in reader could otherwise trigger them.
func TestComicStatsAdminRoutesRequireAuthorization(t *testing.T) {
	cases := []struct {
		method string
		path   string
	}{
		{http.MethodPost, "/admin/comic-stats/trigger/11111111-1111-1111-1111-111111111111"},
		{http.MethodPost, "/admin/comic-stats/trigger-all"},
	}

	for _, tc := range cases {
		t.Run(tc.method+" "+tc.path, func(t *testing.T) {
			chain := handlerChainFor(t, registerComicStatsRoutes, tc.method, tc.path)
			assertGuardedByAuthz(t, chain, tc.method, tc.path)
		})
	}
}

func TestAuthorizationAdminRoutesRequireAuthorization(t *testing.T) {
	cases := []struct {
		name     string
		register func(*gin.Engine)
		method   string
		path     string
	}{
		{"users list", registerUserRoutes, http.MethodGet, "/users"},
		{"user authorization", registerUserRoutes, http.MethodGet, "/users/11111111-1111-1111-1111-111111111111/authorization"},
		{"replace user roles", registerUserRoutes, http.MethodPost, "/users/11111111-1111-1111-1111-111111111111/roles"},
		{"remove user role", registerUserRoutes, http.MethodDelete, "/users/11111111-1111-1111-1111-111111111111/roles/22222222-2222-2222-2222-222222222222"},
		{"roles list", registerRoleRoutes, http.MethodGet, "/roles/all"},
		{"permission catalog", registerPermissionRoutes, http.MethodGet, "/permissions"},
		{"role detail", registerRoleRoutes, http.MethodGet, "/roles/11111111-1111-1111-1111-111111111111"},
		{"create role", registerRoleRoutes, http.MethodPost, "/roles"},
		{"update role", registerRoleRoutes, http.MethodPut, "/roles/11111111-1111-1111-1111-111111111111"},
		{"delete role", registerRoleRoutes, http.MethodDelete, "/roles/11111111-1111-1111-1111-111111111111"},
		{"replace role permissions", registerRoleRoutes, http.MethodPost, "/roles/11111111-1111-1111-1111-111111111111/permissions"},
		{"remove role permission", registerRoleRoutes, http.MethodDelete, "/roles/11111111-1111-1111-1111-111111111111/permissions/role:manage"},
		{"authorization audit", registerAuthorizationRoutes, http.MethodGet, "/authorization/audit-logs"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			chain := handlerChainFor(t, tc.register, tc.method, tc.path)
			assertGuardedByAuthz(t, chain, tc.method, tc.path)
		})
	}
}

func registerUserRoutes(e *gin.Engine) {
	userroute.NewUserRoute(userroute.UserRouteParams{
		R:               e,
		UserHandler:     userroute.NewUserHandler(userroute.UserHandlerParams{}),
		AuthMiddleware:  &authmiddleware.AuthMiddleware{},
		AuthzMiddleware: &authzmiddleware.AuthzMiddleware{},
	}).Setup()
}

func registerRoleRoutes(e *gin.Engine) {
	roleroute.NewRoleRoute(roleroute.RoleRouteParams{
		R:               e,
		RoleHandler:     roleroute.NewRoleHandler(roleroute.RoleHandlerParams{}),
		AuthMiddleware:  &authmiddleware.AuthMiddleware{},
		AuthzMiddleware: &authzmiddleware.AuthzMiddleware{},
	}).Setup()
}

func registerPermissionRoutes(e *gin.Engine) {
	permissionroute.NewPermissionRoute(permissionroute.PermissionRouteParams{
		R:                 e,
		PermissionHandler: permissionroute.NewPermissionHandler(permissionroute.PermissionHandlerParams{}),
		AuthMiddleware:    &authmiddleware.AuthMiddleware{},
		AuthzMiddleware:   &authzmiddleware.AuthzMiddleware{},
	}).Setup()
}

func registerAuthorizationRoutes(e *gin.Engine) {
	authorizationroute.NewRoute(authorizationroute.RouteParams{
		Engine:          e,
		Handler:         authorizationroute.NewHandler(authorizationroute.HandlerParams{}),
		AuthMiddleware:  &authmiddleware.AuthMiddleware{},
		AuthzMiddleware: &authzmiddleware.AuthzMiddleware{},
	}).Setup()
}
