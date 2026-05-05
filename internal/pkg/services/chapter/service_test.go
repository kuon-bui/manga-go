package chapterserivce

import (
	"context"
	"net/http"
	"testing"

	"manga-go/internal/pkg/common"
	"manga-go/internal/pkg/logger"
	"manga-go/internal/pkg/model"
	chapterrepo "manga-go/internal/pkg/repo/chapter"
	comicrepo "manga-go/internal/pkg/repo/comic"
	readingprogressrepo "manga-go/internal/pkg/repo/reading_progress"
	usercomicreadrepo "manga-go/internal/pkg/repo/user_comic_read"
	chapterrequest "manga-go/internal/pkg/request/chapter"
	"manga-go/internal/pkg/testutil"

	"github.com/glebarez/sqlite"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func newChapterService(t *testing.T, createTables bool) *ChapterService {
	t.Helper()

	db, err := gorm.Open(sqlite.Open(testutil.NewSQLiteDSN()), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open sqlite db: %v", err)
	}

	if createTables {
		testutil.MustSyncSchemas(t, db,
			&testutil.Comic{},
			&testutil.Chapter{},
			&testutil.Page{},
			&testutil.ReadingProgress{},
			&testutil.UserComicRead{},
		)
	}

	return &ChapterService{
		logger:              logger.NewLogger(),
		chapterRepo:         chapterrepo.NewChapterRepo(db, nil),
		comicRepo:           comicrepo.NewComicRepo(db),
		readingProgressRepo: readingprogressrepo.NewReadingProgressRepo(db),
		userComicReadRepo:   usercomicreadrepo.NewUserComicReadRepo(db),
	}
}

func TestListChaptersReturnsErrorWithoutComicContext(t *testing.T) {
	t.Parallel()

	s := newChapterService(t, true)
	res := s.ListChapters(context.Background(), &common.Paging{Page: 1, Limit: 10})

	if res.Success {
		t.Fatalf("expected failure result")
	}
	if res.HttpStatus != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", res.HttpStatus)
	}
	if res.Message != "Comic not found in context" {
		t.Fatalf("unexpected message: %s", res.Message)
	}
}

func TestCreateChapterReturnsErrorWithoutComicContext(t *testing.T) {
	t.Parallel()

	s := newChapterService(t, true)
	res := s.CreateChapter(context.Background(), &chapterrequest.CreateChapterRequest{
		Number: "1",
		Title:  "Chapter 1",
		Slug:   "chapter-1",
		Pages: []chapterrequest.PageRequest{{
			PageType: common.ContentTypeImage,
			ImageURL: "https://example.com/1.jpg",
		}},
	})

	if res.Success {
		t.Fatalf("expected failure result")
	}
	if res.HttpStatus != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", res.HttpStatus)
	}
	if res.Message != "Comic not found in context" {
		t.Fatalf("unexpected message: %s", res.Message)
	}
}

func TestGetChapterReturnsNotFoundWhenMissing(t *testing.T) {
	t.Parallel()

	s := newChapterService(t, true)
	ctx := common.SetComicIdToContext(context.Background(), uuid.New())
	res := s.GetChapter(ctx, "missing-chapter")

	if res.Success {
		t.Fatalf("expected failure result")
	}
	if res.HttpStatus != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", res.HttpStatus)
	}
	if res.Message != "Chapter not found" {
		t.Fatalf("unexpected message: %s", res.Message)
	}
}

func TestUpdateChapterReturnsNotFoundWhenMissing(t *testing.T) {
	t.Parallel()

	s := newChapterService(t, true)
	ctx := common.SetComicIdToContext(context.Background(), uuid.New())
	res := s.UpdateChapter(ctx, "missing-chapter", &chapterrequest.UpdateChapterRequest{
		Number: "2",
		Title:  "Updated",
		Slug:   "updated",
	})

	if res.Success {
		t.Fatalf("expected failure result")
	}
	if res.HttpStatus != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", res.HttpStatus)
	}
	if res.Message != "Chapter not found" {
		t.Fatalf("unexpected message: %s", res.Message)
	}
}

func TestPublishChapterReturnsNotFoundWhenMissing(t *testing.T) {
	t.Parallel()

	s := newChapterService(t, true)
	ctx := common.SetComicIdToContext(context.Background(), uuid.New())
	res := s.PublishChapter(ctx, "missing-chapter", &chapterrequest.PublishChapterRequest{IsPublished: true})

	if res.Success {
		t.Fatalf("expected failure result")
	}
	if res.HttpStatus != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", res.HttpStatus)
	}
	if res.Message != "Chapter not found" {
		t.Fatalf("unexpected message: %s", res.Message)
	}
}

func TestUpdateReadingProgressReturnsNotFoundWhenChapterMissing(t *testing.T) {
	t.Parallel()

	s := newChapterService(t, true)
	res := s.UpdateReadingProgress(context.Background(), &model.User{SqlModel: common.SqlModel{ID: uuid.New()}}, uuid.New(), &chapterrequest.UpdateReadingProgressRequest{ScrollPercent: 30})

	if res.Success {
		t.Fatalf("expected failure result")
	}
	if res.HttpStatus != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", res.HttpStatus)
	}
	if res.Message != "Chapter not found" {
		t.Fatalf("unexpected message: %s", res.Message)
	}
}

func TestCreateChapterReturnsDbErrorWhenTableMissing(t *testing.T) {
	t.Parallel()

	s := newChapterService(t, false)
	ctx := common.SetComicIdToContext(context.Background(), uuid.New())
	res := s.CreateChapter(ctx, &chapterrequest.CreateChapterRequest{
		Number: "1",
		Title:  "Chapter 1",
		Slug:   "chapter-1",
		Pages: []chapterrequest.PageRequest{{
			PageType: common.ContentTypeImage,
			ImageURL: "https://example.com/1.jpg",
		}},
	})

	if res.Success {
		t.Fatalf("expected failure result")
	}
	if res.HttpStatus != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", res.HttpStatus)
	}
	if res.Message != "database error" {
		t.Fatalf("unexpected message: %s", res.Message)
	}
	if res.Error == nil {
		t.Fatalf("expected non-nil error")
	}
}
