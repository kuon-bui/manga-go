package services_test

import (
	"testing"

	"manga-go/internal/pkg/services"
	authorservice "manga-go/internal/pkg/services/author"
	chapterservice "manga-go/internal/pkg/services/chapter"
	comicservice "manga-go/internal/pkg/services/comic"
	commentservice "manga-go/internal/pkg/services/comment"
	fileservice "manga-go/internal/pkg/services/file"
	genreservice "manga-go/internal/pkg/services/genre"
	permissionservice "manga-go/internal/pkg/services/permission"
	ratingservice "manga-go/internal/pkg/services/rating"
	readinghistoryservice "manga-go/internal/pkg/services/reading_history"
	roleservice "manga-go/internal/pkg/services/role"
	tagservice "manga-go/internal/pkg/services/tag"
	translationgroupservice "manga-go/internal/pkg/services/translation_group"
	userservice "manga-go/internal/pkg/services/user"
)

func assertNotNil(t *testing.T, name string, value any) {
	t.Helper()
	if value == nil {
		t.Fatalf("expected %s to be non-nil", name)
	}
}

func TestServicesModulesAreRegistered(t *testing.T) {
	assertNotNil(t, "services.Module", services.Module)
	assertNotNil(t, "userservice.Module", userservice.Module)
	assertNotNil(t, "authorservice.Module", authorservice.Module)
	assertNotNil(t, "genreservice.Module", genreservice.Module)
	assertNotNil(t, "fileservice.Module", fileservice.Module)
	assertNotNil(t, "tagservice.Module", tagservice.Module)
	assertNotNil(t, "comicservice.Module", comicservice.Module)
	assertNotNil(t, "chapterservice.Module", chapterservice.Module)
	assertNotNil(t, "translationgroupservice.Module", translationgroupservice.Module)
	assertNotNil(t, "roleservice.Module", roleservice.Module)
	assertNotNil(t, "permissionservice.Module", permissionservice.Module)
	assertNotNil(t, "ratingservice.Module", ratingservice.Module)
	assertNotNil(t, "readinghistoryservice.Module", readinghistoryservice.Module)
	assertNotNil(t, "commentservice.Module", commentservice.Module)
}

func TestServiceConstructorsReturnInstance(t *testing.T) {
	if s := userservice.NewUserService(userservice.UserServiceParams{}); s == nil {
		t.Fatal("expected user service instance")
	}
	if s := authorservice.NewAuthorService(authorservice.AuthorServiceParams{}); s == nil {
		t.Fatal("expected author service instance")
	}
	if s := genreservice.NewGenreService(genreservice.GenreServiceParams{}); s == nil {
		t.Fatal("expected genre service instance")
	}
	if s := fileservice.NewFileService(fileservice.FileServiceParams{}); s == nil {
		t.Fatal("expected file service instance")
	}
	if s := tagservice.NewTagService(tagservice.TagServiceParams{}); s == nil {
		t.Fatal("expected tag service instance")
	}
	if s := comicservice.NewComicService(comicservice.ComicServiceParams{}); s == nil {
		t.Fatal("expected comic service instance")
	}
	if s := chapterservice.NewChapterService(chapterservice.ChapterServiceParams{}); s == nil {
		t.Fatal("expected chapter service instance")
	}
	if s := translationgroupservice.NewTranslationGroupService(translationgroupservice.TranslationGroupServiceParams{}); s == nil {
		t.Fatal("expected translation group service instance")
	}
	if s := roleservice.NewRoleService(roleservice.RoleServiceParams{}); s == nil {
		t.Fatal("expected role service instance")
	}
	if s := permissionservice.NewPermissionService(permissionservice.PermissionServiceParams{}); s == nil {
		t.Fatal("expected permission service instance")
	}
	if s := ratingservice.NewRatingService(ratingservice.RatingServiceParams{}); s == nil {
		t.Fatal("expected rating service instance")
	}
	if s := readinghistoryservice.NewReadingHistoryService(readinghistoryservice.ReadingHistoryServiceParams{}); s == nil {
		t.Fatal("expected reading history service instance")
	}
	if s := commentservice.NewCommentService(commentservice.CommentServiceParams{}); s == nil {
		t.Fatal("expected comment service instance")
	}
}
