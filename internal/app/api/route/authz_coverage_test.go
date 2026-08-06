package route_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	comicstatsroute "manga-go/internal/app/api/route/comic_stats"
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
