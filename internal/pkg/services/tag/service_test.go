package tagservice

import (
	"context"
	"net/http"
	"reflect"
	"testing"

	"manga-go/internal/pkg/common"
	"manga-go/internal/pkg/logger"
	tagrepo "manga-go/internal/pkg/repo/tag"
	tagrequest "manga-go/internal/pkg/request/tag"
	"manga-go/internal/pkg/testutil"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func newTagService(t *testing.T, createTable bool) *TagService {
	t.Helper()

	db, err := gorm.Open(sqlite.Open(testutil.NewSQLiteDSN()), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open sqlite db: %v", err)
	}

	if createTable {
		testutil.MustSyncSchemas(t, db, &testutil.Tag{})
	}

	return &TagService{
		logger:  logger.NewLogger(),
		tagRepo: tagrepo.NewTagRepo(db, nil),
	}
}

func paginationTotalTag(data any) int64 {
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

func TestListTagsReturnsEmptyPagination(t *testing.T) {
	t.Parallel()

	s := newTagService(t, true)
	res := s.ListTags(context.Background(), &common.Paging{Page: 1, Limit: 10})

	if !res.Success {
		t.Fatalf("expected success result")
	}
	if res.HttpStatus != http.StatusOK {
		t.Fatalf("expected status 200, got %d", res.HttpStatus)
	}
	if res.Message != "Tags retrieved successfully" {
		t.Fatalf("unexpected message: %s", res.Message)
	}
	if total := paginationTotalTag(res.Data); total != 0 {
		t.Fatalf("expected total 0, got %d", total)
	}
}

func TestDeleteTagReturnsNotFoundWhenMissing(t *testing.T) {
	t.Parallel()

	s := newTagService(t, true)
	res := s.DeleteTag(context.Background(), "missing-tag")

	if res.Success {
		t.Fatalf("expected failure result")
	}
	if res.HttpStatus != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", res.HttpStatus)
	}
	if res.Message != "Tag not found" {
		t.Fatalf("unexpected message: %s", res.Message)
	}
}

func TestCreateTagReturnsDbErrorWhenTableMissing(t *testing.T) {
	t.Parallel()

	s := newTagService(t, false)
	res := s.CreateTag(context.Background(), &tagrequest.CreateTagRequest{
		Name: "Shounen",
		Slug: "shounen",
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

func TestDeleteTagReturnsDbErrorWhenTableMissing(t *testing.T) {
	t.Parallel()

	s := newTagService(t, false)
	res := s.DeleteTag(context.Background(), "missing-tag")

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
