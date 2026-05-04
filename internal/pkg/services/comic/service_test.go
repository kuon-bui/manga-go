package comicservice

import (
	"context"
	"net/http"
	"reflect"
	"testing"

	"manga-go/internal/pkg/common"
	"manga-go/internal/pkg/constant"
	"manga-go/internal/pkg/logger"
	comicrepo "manga-go/internal/pkg/repo/comic"
	genrerepo "manga-go/internal/pkg/repo/genre"
	tagrepo "manga-go/internal/pkg/repo/tag"
	usercomicreadrepo "manga-go/internal/pkg/repo/user_comic_read"
	comicrequest "manga-go/internal/pkg/request/comic"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func newComicService(t *testing.T, createTable bool) *ComicService {
	t.Helper()

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open sqlite db: %v", err)
	}

	if createTable {
		err = db.Exec(`
			CREATE TABLE comics (
				id TEXT PRIMARY KEY,
				title TEXT,
				slug TEXT,
				alternative_titles TEXT,
				description TEXT,
				thumbnail TEXT,
				banner TEXT,
				type TEXT,
				status TEXT,
				age_rating TEXT,
				is_published BOOLEAN,
				is_hot BOOLEAN,
				is_featured BOOLEAN,
				published_year INTEGER,
				last_chapter_at DATETIME,
				artist_id TEXT,
				created_at DATETIME,
				updated_at DATETIME,
				deleted_at DATETIME
			)
		`).Error
		if err != nil {
			t.Fatalf("failed to create comics table: %v", err)
		}
	}

	return &ComicService{
		logger:            logger.NewLogger(),
		comicRepo:         comicrepo.NewComicRepo(db),
		genreRepo:         genrerepo.NewGenreRepo(db),
		tagRepo:           tagrepo.NewTagRepo(db, nil),
		userComicReadRepo: usercomicreadrepo.NewUserComicReadRepo(db),
	}
}

func comicPaginationTotal(data any) int64 {
	v := reflect.ValueOf(data)
	if !v.IsValid() {
		return -1
	}

	field := v.FieldByName("Total")
	if !field.IsValid() || field.Kind() != reflect.Int64 {
		return -1
	}

	return field.Int()
}

func TestListComicsReturnsEmptyPagination(t *testing.T) {
	t.Parallel()

	s := newComicService(t, true)
	res := s.ListComics(context.Background(), &common.Paging{Page: 1, Limit: 10})

	if !res.Success {
		t.Fatalf("expected success result")
	}
	if res.HttpStatus != http.StatusOK {
		t.Fatalf("expected status 200, got %d", res.HttpStatus)
	}
	if res.Message != "Comics retrieved successfully" {
		t.Fatalf("unexpected message: %s", res.Message)
	}
	if total := comicPaginationTotal(res.Data); total != 0 {
		t.Fatalf("expected total 0, got %d", total)
	}
}

func TestGetComicReturnsNotFoundWhenMissing(t *testing.T) {
	t.Parallel()

	s := newComicService(t, true)
	res := s.GetComic(context.Background(), "missing-comic")

	if res.Success {
		t.Fatalf("expected failure result")
	}
	if res.HttpStatus != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", res.HttpStatus)
	}
	if res.Message != "Comic not found" {
		t.Fatalf("unexpected message: %s", res.Message)
	}
}

func TestUpdateComicReturnsNotFoundWhenMissing(t *testing.T) {
	t.Parallel()

	s := newComicService(t, true)
	res := s.UpdateComic(context.Background(), "missing-comic", &comicrequest.UpdateComicRequest{
		Title: "Updated title",
		Slug:  "updated-slug",
		Type:  constant.ComicTypeManga,
	})

	if res.Success {
		t.Fatalf("expected failure result")
	}
	if res.HttpStatus != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", res.HttpStatus)
	}
	if res.Message != "Comic not found" {
		t.Fatalf("unexpected message: %s", res.Message)
	}
}

func TestDeleteComicReturnsNotFoundWhenMissing(t *testing.T) {
	t.Parallel()

	s := newComicService(t, true)
	res := s.DeleteComic(context.Background(), "missing-comic")

	if res.Success {
		t.Fatalf("expected failure result")
	}
	if res.HttpStatus != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", res.HttpStatus)
	}
	if res.Message != "Comic not found" {
		t.Fatalf("unexpected message: %s", res.Message)
	}
}

func TestUpdateComicStatusReturnsNotFoundWhenMissing(t *testing.T) {
	t.Parallel()

	s := newComicService(t, true)
	res := s.UpdateComicStatus(context.Background(), "missing-comic", &comicrequest.UpdateComicStatusRequest{
		Status: constant.ComicStatusCompleted,
	})

	if res.Success {
		t.Fatalf("expected failure result")
	}
	if res.HttpStatus != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", res.HttpStatus)
	}
	if res.Message != "Comic not found" {
		t.Fatalf("unexpected message: %s", res.Message)
	}
}

func TestPublishComicReturnsNotFoundWhenMissing(t *testing.T) {
	t.Parallel()

	s := newComicService(t, true)
	res := s.PublishComic(context.Background(), "missing-comic", &comicrequest.PublishComicRequest{IsPublished: true})

	if res.Success {
		t.Fatalf("expected failure result")
	}
	if res.HttpStatus != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", res.HttpStatus)
	}
	if res.Message != "Comic not found" {
		t.Fatalf("unexpected message: %s", res.Message)
	}
}

func TestCreateComicReturnsDbErrorWhenTableMissing(t *testing.T) {
	t.Parallel()

	s := newComicService(t, false)
	res := s.CreateComic(context.Background(), &comicrequest.CreateComicRequest{
		Title: "One Piece",
		Slug:  "one-piece",
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
