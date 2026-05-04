package route_test

import (
	"testing"

	"manga-go/internal/app/api/route"
	authorroute "manga-go/internal/app/api/route/author"
	chapterroute "manga-go/internal/app/api/route/chapter"
	comicroute "manga-go/internal/app/api/route/comic"
	commentroute "manga-go/internal/app/api/route/comment"
	fileroute "manga-go/internal/app/api/route/file"
	genreroute "manga-go/internal/app/api/route/genre"
	permissionroute "manga-go/internal/app/api/route/permission"
	ratingroute "manga-go/internal/app/api/route/rating"
	readinghistoryroute "manga-go/internal/app/api/route/reading_history"
	roleroute "manga-go/internal/app/api/route/role"
	swaggerrouter "manga-go/internal/app/api/route/swagger"
	tagroute "manga-go/internal/app/api/route/tag"
	translationgrouproute "manga-go/internal/app/api/route/translation_group"
	userroute "manga-go/internal/app/api/route/user"
	authmiddleware "manga-go/internal/app/middleware/auth"
	slugmiddleware "manga-go/internal/app/middleware/slug"

	"github.com/gin-gonic/gin"
)

func assertNotNil(t *testing.T, name string, value any) {
	t.Helper()
	if value == nil {
		t.Fatalf("expected %s to be non-nil", name)
	}
}

func TestRouteModulesAreRegistered(t *testing.T) {
	assertNotNil(t, "route.Module", route.Module)
	assertNotNil(t, "userroute.Module", userroute.Module)
	assertNotNil(t, "authorroute.Module", authorroute.Module)
	assertNotNil(t, "genreroute.Module", genreroute.Module)
	assertNotNil(t, "fileroute.Module", fileroute.Module)
	assertNotNil(t, "tagroute.Module", tagroute.Module)
	assertNotNil(t, "comicroute.Module", comicroute.Module)
	assertNotNil(t, "chapterroute.Module", chapterroute.Module)
	assertNotNil(t, "translationgrouproute.Module", translationgrouproute.Module)
	assertNotNil(t, "roleroute.Module", roleroute.Module)
	assertNotNil(t, "permissionroute.Module", permissionroute.Module)
	assertNotNil(t, "ratingroute.Module", ratingroute.Module)
	assertNotNil(t, "readinghistoryroute.Module", readinghistoryroute.Module)
	assertNotNil(t, "commentroute.Module", commentroute.Module)
	assertNotNil(t, "swaggerrouter.Module", swaggerrouter.Module)
}

func assertRouteExists(t *testing.T, routes gin.RoutesInfo, method string, path string) {
	t.Helper()
	for _, route := range routes {
		if route.Method == method && route.Path == path {
			return
		}
	}
	t.Fatalf("expected route %s %s to be registered", method, path)
}

func TestRouteSetupRegistersEndpoints(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("author", func(t *testing.T) {
		e := gin.New()
		r := authorroute.NewAuthorRoute(authorroute.AuthorRouteParams{R: e, AuthorHandler: authorroute.NewAuthorHandler(authorroute.AuthorHandlerParams{}), AuthMiddleware: &authmiddleware.AuthMiddleware{}})
		r.Setup()

		routes := e.Routes()
		if len(routes) != 6 {
			t.Fatalf("expected 6 routes, got %d", len(routes))
		}
		assertRouteExists(t, routes, "GET", "/authors/")
	})

	t.Run("genre", func(t *testing.T) {
		e := gin.New()
		r := genreroute.NewGenreRoute(genreroute.GenreRouteParams{R: e, GenreHandler: genreroute.NewGenreHandler(genreroute.GenreHandlerParams{}), AuthMiddleware: &authmiddleware.AuthMiddleware{}})
		r.Setup()

		routes := e.Routes()
		if len(routes) != 6 {
			t.Fatalf("expected 6 routes, got %d", len(routes))
		}
		assertRouteExists(t, routes, "GET", "/genres")
	})

	t.Run("tag", func(t *testing.T) {
		e := gin.New()
		r := tagroute.NewTagRoute(tagroute.TagRouteParams{R: e, TagHandler: tagroute.NewTagHandler(tagroute.TagHandlerParams{}), AuthMiddleware: &authmiddleware.AuthMiddleware{}})
		r.Setup()

		routes := e.Routes()
		if len(routes) != 4 {
			t.Fatalf("expected 4 routes, got %d", len(routes))
		}
		assertRouteExists(t, routes, "GET", "/tags")
	})

	t.Run("comic", func(t *testing.T) {
		e := gin.New()
		r := comicroute.NewComicRoute(comicroute.ComicRouteParams{R: e, ComicHandler: comicroute.NewComicHandler(comicroute.ComicHandlerParams{}), AuthMiddleware: &authmiddleware.AuthMiddleware{}})
		r.Setup()

		routes := e.Routes()
		if len(routes) != 7 {
			t.Fatalf("expected 7 routes, got %d", len(routes))
		}
		assertRouteExists(t, routes, "GET", "/comics")
	})

	t.Run("chapter", func(t *testing.T) {
		e := gin.New()
		r := chapterroute.NewChapterRoute(chapterroute.ChapterRouteParams{R: e, Handler: chapterroute.NewChapterHandler(chapterroute.ChapterHandlerParams{}), AuthMiddleware: &authmiddleware.AuthMiddleware{}, SlugMiddleware: &slugmiddleware.SlugMiddleware{}})
		r.Setup()

		routes := e.Routes()
		if len(routes) != 9 {
			t.Fatalf("expected 9 routes, got %d", len(routes))
		}
		assertRouteExists(t, routes, "GET", "/comics/:comicSlug/chapters")
	})

	t.Run("translation-group", func(t *testing.T) {
		e := gin.New()
		r := translationgrouproute.NewTranslationGroupRoute(translationgrouproute.TranslationGroupRouteParams{R: e, TranslationGroupHandler: translationgrouproute.NewTranslationGroupHandler(translationgrouproute.TranslationGroupHandlerParams{}), AuthMiddleware: &authmiddleware.AuthMiddleware{}})
		r.Setup()

		routes := e.Routes()
		if len(routes) != 6 {
			t.Fatalf("expected 6 routes, got %d", len(routes))
		}
		assertRouteExists(t, routes, "GET", "/translation-groups/")
	})

	t.Run("role", func(t *testing.T) {
		e := gin.New()
		r := roleroute.NewRoleRoute(roleroute.RoleRouteParams{R: e, RoleHandler: roleroute.NewRoleHandler(roleroute.RoleHandlerParams{}), AuthMiddleware: &authmiddleware.AuthMiddleware{}})
		r.Setup()

		routes := e.Routes()
		if len(routes) != 8 {
			t.Fatalf("expected 8 routes, got %d", len(routes))
		}
		assertRouteExists(t, routes, "GET", "/roles/")
	})

	t.Run("permission", func(t *testing.T) {
		e := gin.New()
		r := permissionroute.NewPermissionRoute(permissionroute.PermissionRouteParams{R: e, PermissionHandler: permissionroute.NewPermissionHandler(permissionroute.PermissionHandlerParams{}), AuthMiddleware: &authmiddleware.AuthMiddleware{}})
		r.Setup()

		routes := e.Routes()
		if len(routes) != 5 {
			t.Fatalf("expected 5 routes, got %d", len(routes))
		}
		assertRouteExists(t, routes, "GET", "/permissions/")
	})

	t.Run("rating", func(t *testing.T) {
		e := gin.New()
		r := ratingroute.NewRatingRoute(ratingroute.RatingRouteParams{R: e, RatingHandler: ratingroute.NewRatingHandler(ratingroute.RatingHandlerParams{}), AuthMiddleware: &authmiddleware.AuthMiddleware{}, SlugMiddleware: &slugmiddleware.SlugMiddleware{}})
		r.Setup()

		routes := e.Routes()
		if len(routes) != 5 {
			t.Fatalf("expected 5 routes, got %d", len(routes))
		}
		assertRouteExists(t, routes, "GET", "/ratings/comics/:comicSlug")
	})

	t.Run("reading-history", func(t *testing.T) {
		e := gin.New()
		r := readinghistoryroute.NewReadingHistoryRoute(readinghistoryroute.ReadingHistoryRouteParams{R: e, Handler: readinghistoryroute.NewReadingHistoryHandler(readinghistoryroute.ReadingHistoryHandlerParams{}), AuthMiddleware: &authmiddleware.AuthMiddleware{}})
		r.Setup()

		routes := e.Routes()
		if len(routes) != 5 {
			t.Fatalf("expected 5 routes, got %d", len(routes))
		}
		assertRouteExists(t, routes, "GET", "/reading-histories")
	})

	t.Run("comment", func(t *testing.T) {
		e := gin.New()
		r := commentroute.NewCommentRoute(commentroute.CommentRouteParams{R: e, CommentHandler: commentroute.NewCommentHandler(commentroute.CommentHandlerParams{}), AuthMiddleware: &authmiddleware.AuthMiddleware{}})
		r.Setup()

		routes := e.Routes()
		if len(routes) != 6 {
			t.Fatalf("expected 6 routes, got %d", len(routes))
		}
		assertRouteExists(t, routes, "GET", "/comments")
	})

	t.Run("file", func(t *testing.T) {
		e := gin.New()
		r := fileroute.NewFileRoute(fileroute.FileRouteParams{R: e, FileHandler: fileroute.NewFileHandler(fileroute.FileHandlerParams{}), AuthMiddleware: &authmiddleware.AuthMiddleware{}})
		r.Setup()

		routes := e.Routes()
		if len(routes) != 3 {
			t.Fatalf("expected 3 routes, got %d", len(routes))
		}
		assertRouteExists(t, routes, "POST", "/files/upload")
	})

	t.Run("user", func(t *testing.T) {
		e := gin.New()
		r := userroute.NewUserRoute(userroute.UserRouteParams{R: e, UserHandler: userroute.NewUserHandler(userroute.UserHandlerParams{}), AuthMiddleware: &authmiddleware.AuthMiddleware{}})
		r.Setup()

		routes := e.Routes()
		if len(routes) != 10 {
			t.Fatalf("expected 10 routes, got %d", len(routes))
		}
		assertRouteExists(t, routes, "POST", "/users/")
	})

	t.Run("swagger", func(t *testing.T) {
		e := gin.New()
		r := swaggerrouter.NewSwaggerRoute(swaggerrouter.SwaggerRouteParams{R: e, AuthMiddleware: &authmiddleware.AuthMiddleware{}})
		r.Setup()

		routes := e.Routes()
		if len(routes) != 1 {
			t.Fatalf("expected 1 route, got %d", len(routes))
		}
		assertRouteExists(t, routes, "GET", "/swagger/*any")
	})
}
