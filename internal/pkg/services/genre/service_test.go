package genreservice

import (
	"context"
	"net/http"
	"reflect"
	"testing"

	"manga-go/internal/pkg/common"
	"manga-go/internal/pkg/logger"
	genrerepo "manga-go/internal/pkg/repo/genre"
	genrerequest "manga-go/internal/pkg/request/genre"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func newGenreService(t *testing.T, createTable bool) *GenreService {
	t.Helper()

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open sqlite db: %v", err)
	}

	if createTable {
		err = db.Exec(`
			CREATE TABLE genres (
				id TEXT PRIMARY KEY,
				name TEXT,
				slug TEXT,
				description TEXT,
				thumbnail TEXT,
				created_at DATETIME,
				updated_at DATETIME,
				deleted_at DATETIME
			)
		`).Error
		if err != nil {
			t.Fatalf("failed to create genres table: %v", err)
		}
	}

	return &GenreService{
		logger:    logger.NewLogger(),
		genreRepo: genrerepo.NewGenreRepo(db),
	}
}

func paginationTotal(data any) int64 {
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

func TestGetGenreReturnsNotFoundWhenRecordMissing(t *testing.T) {
	t.Parallel()

	s := newGenreService(t, true)
	res := s.GetGenre(context.Background(), "missing-genre")

	if res.Success {
		t.Fatalf("expected failure result")
	}
	if res.HttpStatus != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", res.HttpStatus)
	}
	if res.Message != "Genre not found" {
		t.Fatalf("unexpected message: %s", res.Message)
	}
}

func TestGetGenreReturnsDbErrorWhenTableMissing(t *testing.T) {
	t.Parallel()

	s := newGenreService(t, false)
	res := s.GetGenre(context.Background(), "missing-genre")

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

func TestListGenresReturnsEmptyPagination(t *testing.T) {
	t.Parallel()

	s := newGenreService(t, true)
	res := s.ListGenres(context.Background(), &common.Paging{Page: 1, Limit: 10})

	if !res.Success {
		t.Fatalf("expected success result")
	}
	if res.HttpStatus != http.StatusOK {
		t.Fatalf("expected status 200, got %d", res.HttpStatus)
	}
	if res.Message != "Genres retrieved successfully" {
		t.Fatalf("unexpected message: %s", res.Message)
	}
	if total := paginationTotal(res.Data); total != 0 {
		t.Fatalf("expected total 0, got %d", total)
	}
}

func TestCreateGenreReturnsDbErrorWhenTableMissing(t *testing.T) {
	t.Parallel()

	s := newGenreService(t, false)
	res := s.CreateGenre(context.Background(), &genrerequest.CreateGenreRequest{
		Name: "Action",
		Slug: "action",
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
