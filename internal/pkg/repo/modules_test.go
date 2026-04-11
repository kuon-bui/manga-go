package repo_test

import (
	"testing"

	"manga-go/internal/pkg/repo"
	authorrepo "manga-go/internal/pkg/repo/author"
	chapterrepo "manga-go/internal/pkg/repo/chapter"
	comicrepo "manga-go/internal/pkg/repo/comic"
	commentrepo "manga-go/internal/pkg/repo/comment"
	genrerepo "manga-go/internal/pkg/repo/genre"
	permissionrepo "manga-go/internal/pkg/repo/permission"
	ratingrepo "manga-go/internal/pkg/repo/rating"
	reactionrepo "manga-go/internal/pkg/repo/reaction"
	readinghistoryrepo "manga-go/internal/pkg/repo/reading_history"
	readingprogressrepo "manga-go/internal/pkg/repo/reading_progress"
	rolerepo "manga-go/internal/pkg/repo/role"
	tagrepo "manga-go/internal/pkg/repo/tag"
	translationgrouprepo "manga-go/internal/pkg/repo/translation_group"
	userrepo "manga-go/internal/pkg/repo/user"
	usercomicreadrepo "manga-go/internal/pkg/repo/user_comic_read"
)

func assertNotNil(t *testing.T, name string, value any) {
	t.Helper()
	if value == nil {
		t.Fatalf("expected %s to be non-nil", name)
	}
}

func TestRepoModulesAreRegistered(t *testing.T) {
	assertNotNil(t, "repo.Module", repo.Module)
	assertNotNil(t, "userrepo.Module", userrepo.Module)
	assertNotNil(t, "authorrepo.Module", authorrepo.Module)
	assertNotNil(t, "genrerepo.Module", genrerepo.Module)
	assertNotNil(t, "tagrepo.Module", tagrepo.Module)
	assertNotNil(t, "comicrepo.Module", comicrepo.Module)
	assertNotNil(t, "chapterrepo.Module", chapterrepo.Module)
	assertNotNil(t, "translationgrouprepo.Module", translationgrouprepo.Module)
	assertNotNil(t, "rolerepo.Module", rolerepo.Module)
	assertNotNil(t, "permissionrepo.Module", permissionrepo.Module)
	assertNotNil(t, "ratingrepo.Module", ratingrepo.Module)
	assertNotNil(t, "readinghistoryrepo.Module", readinghistoryrepo.Module)
	assertNotNil(t, "commentrepo.Module", commentrepo.Module)
	assertNotNil(t, "reactionrepo.Module", reactionrepo.Module)
	assertNotNil(t, "readingprogressrepo.Module", readingprogressrepo.Module)
	assertNotNil(t, "usercomicreadrepo.Module", usercomicreadrepo.Module)
}

func TestRepoConstructorsReturnInstance(t *testing.T) {
	if r := userrepo.NewUserRepository(nil, nil); r == nil {
		t.Fatal("expected user repo instance")
	}
	if r := authorrepo.NewAuthorRepo(nil); r == nil {
		t.Fatal("expected author repo instance")
	}
	if r := genrerepo.NewGenreRepo(nil); r == nil {
		t.Fatal("expected genre repo instance")
	}
	if r := tagrepo.NewTagRepo(nil, nil); r == nil {
		t.Fatal("expected tag repo instance")
	}
	if r := comicrepo.NewComicRepo(nil); r == nil {
		t.Fatal("expected comic repo instance")
	}
	if r := chapterrepo.NewChapterRepo(nil, nil); r == nil {
		t.Fatal("expected chapter repo instance")
	}
	if r := translationgrouprepo.NewTranslationGroupRepo(nil, nil); r == nil {
		t.Fatal("expected translation group repo instance")
	}
	if r := rolerepo.NewRoleRepo(nil); r == nil {
		t.Fatal("expected role repo instance")
	}
	if r := permissionrepo.NewPermissionRepo(nil); r == nil {
		t.Fatal("expected permission repo instance")
	}
	if r := ratingrepo.NewRatingRepo(nil); r == nil {
		t.Fatal("expected rating repo instance")
	}
	if r := readinghistoryrepo.NewReadingHistoryRepo(nil); r == nil {
		t.Fatal("expected reading history repo instance")
	}
	if r := commentrepo.NewCommentRepo(nil); r == nil {
		t.Fatal("expected comment repo instance")
	}
	if r := reactionrepo.NewReactionRepo(nil); r == nil {
		t.Fatal("expected reaction repo instance")
	}
	if r := readingprogressrepo.NewReadingProgressRepo(nil); r == nil {
		t.Fatal("expected reading progress repo instance")
	}
	if r := usercomicreadrepo.NewUserComicReadRepo(nil); r == nil {
		t.Fatal("expected user comic read repo instance")
	}
}
