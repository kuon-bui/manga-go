package readinghistoryservice

import (
	"context"
	"net/http"
	"reflect"
	"testing"

	"manga-go/internal/pkg/common"
	"manga-go/internal/pkg/logger"
	readinghistoryrepo "manga-go/internal/pkg/repo/reading_history"
	readinghistoryrequest "manga-go/internal/pkg/request/reading_history"

	"github.com/glebarez/sqlite"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func newReadingHistoryService(t *testing.T, createTable bool) *ReadingHistoryService {
	t.Helper()

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open sqlite db: %v", err)
	}

	if createTable {
		err = db.Exec(`
			CREATE TABLE reading_histories (
				id TEXT PRIMARY KEY,
				user_id TEXT,
				chapter_id TEXT,
				comic_id TEXT,
				last_read_at DATETIME,
				created_at DATETIME,
				updated_at DATETIME,
				deleted_at DATETIME
			)
		`).Error
		if err != nil {
			t.Fatalf("failed to create reading_histories table: %v", err)
		}
	}

	return &ReadingHistoryService{
		logger:             logger.NewLogger(),
		readingHistoryRepo: readinghistoryrepo.NewReadingHistoryRepo(db),
	}
}

func readingHistoryPaginationTotal(data any) int64 {
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

func TestListReadingHistoriesReturnsEmptyPagination(t *testing.T) {
	t.Parallel()

	s := newReadingHistoryService(t, true)
	res := s.ListReadingHistories(context.Background(), uuid.New(), &common.Paging{Page: 1, Limit: 10})

	if !res.Success {
		t.Fatalf("expected success result")
	}
	if res.HttpStatus != http.StatusOK {
		t.Fatalf("expected status 200, got %d", res.HttpStatus)
	}
	if res.Message != "Reading histories retrieved successfully" {
		t.Fatalf("unexpected message: %s", res.Message)
	}
	if total := readingHistoryPaginationTotal(res.Data); total != 0 {
		t.Fatalf("expected total 0, got %d", total)
	}
}

func TestGetReadingHistoryReturnsNotFoundWhenMissing(t *testing.T) {
	t.Parallel()

	s := newReadingHistoryService(t, true)
	res := s.GetReadingHistory(context.Background(), uuid.New())

	if res.Success {
		t.Fatalf("expected failure result")
	}
	if res.HttpStatus != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", res.HttpStatus)
	}
	if res.Message != "ReadingHistory not found" {
		t.Fatalf("unexpected message: %s", res.Message)
	}
}

func TestUpdateReadingHistoryReturnsNotFoundWhenMissing(t *testing.T) {
	t.Parallel()

	s := newReadingHistoryService(t, true)
	res := s.UpdateReadingHistory(context.Background(), uuid.New(), &readinghistoryrequest.UpdateReadingHistoryRequest{})

	if res.Success {
		t.Fatalf("expected failure result")
	}
	if res.HttpStatus != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", res.HttpStatus)
	}
	if res.Message != "ReadingHistory not found" {
		t.Fatalf("unexpected message: %s", res.Message)
	}
}

func TestDeleteReadingHistoryReturnsNotFoundWhenMissing(t *testing.T) {
	t.Parallel()

	s := newReadingHistoryService(t, true)
	res := s.DeleteReadingHistory(context.Background(), uuid.New())

	if res.Success {
		t.Fatalf("expected failure result")
	}
	if res.HttpStatus != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", res.HttpStatus)
	}
	if res.Message != "ReadingHistory not found" {
		t.Fatalf("unexpected message: %s", res.Message)
	}
}

func TestCreateReadingHistoryReturnsDbErrorWhenTableMissing(t *testing.T) {
	t.Parallel()

	s := newReadingHistoryService(t, false)
	res := s.CreateReadingHistory(context.Background(), uuid.New(), &readinghistoryrequest.CreateReadingHistoryRequest{
		ChapterID: uuid.New(),
		ComicID:   uuid.New(),
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
