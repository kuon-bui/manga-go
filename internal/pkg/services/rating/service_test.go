package ratingservice

import (
	"context"
	"net/http"
	"reflect"
	"testing"

	"manga-go/internal/pkg/common"
	"manga-go/internal/pkg/logger"
	comicrepo "manga-go/internal/pkg/repo/comic"
	ratingrepo "manga-go/internal/pkg/repo/rating"
	ratingrequest "manga-go/internal/pkg/request/rating"
	"manga-go/internal/pkg/testutil"

	"github.com/glebarez/sqlite"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func newRatingService(t *testing.T, createTable bool) *RatingService {
	t.Helper()

	db, err := gorm.Open(sqlite.Open(testutil.NewSQLiteDSN()), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open sqlite db: %v", err)
	}

	if createTable {
		testutil.MustSyncSchemas(t, db, &testutil.Rating{})
	}

	return &RatingService{
		logger:     logger.NewLogger(),
		ratingRepo: ratingrepo.NewRatingRepo(db),
		comicRepo:  comicrepo.NewComicRepo(db),
	}
}

func ratingPaginationTotal(data any) int64 {
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

func TestListRatingsReturnsEmptyPagination(t *testing.T) {
	t.Parallel()

	s := newRatingService(t, true)
	res := s.ListRatings(context.Background(), uuid.New(), &common.Paging{Page: 1, Limit: 10})

	if !res.Success {
		t.Fatalf("expected success result")
	}
	if res.HttpStatus != http.StatusOK {
		t.Fatalf("expected status 200, got %d", res.HttpStatus)
	}
	if res.Message != "Ratings retrieved successfully" {
		t.Fatalf("unexpected message: %s", res.Message)
	}
	if total := ratingPaginationTotal(res.Data); total != 0 {
		t.Fatalf("expected total 0, got %d", total)
	}
}

func TestUpdateRatingReturnsNotFoundWhenMissing(t *testing.T) {
	t.Parallel()

	s := newRatingService(t, true)
	res := s.UpdateRating(context.Background(), uuid.New(), uuid.New(), &ratingrequest.UpdateRatingRequest{Score: 4})

	if res.Success {
		t.Fatalf("expected failure result")
	}
	if res.HttpStatus != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", res.HttpStatus)
	}
	if res.Message != "Rating not found" {
		t.Fatalf("unexpected message: %s", res.Message)
	}
}

func TestDeleteRatingReturnsNotFoundWhenMissing(t *testing.T) {
	t.Parallel()

	s := newRatingService(t, true)
	res := s.DeleteRating(context.Background(), uuid.New(), uuid.New())

	if res.Success {
		t.Fatalf("expected failure result")
	}
	if res.HttpStatus != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", res.HttpStatus)
	}
	if res.Message != "Rating not found" {
		t.Fatalf("unexpected message: %s", res.Message)
	}
}

func TestCreateRatingReturnsDbErrorWhenTableMissing(t *testing.T) {
	t.Parallel()

	s := newRatingService(t, false)
	res := s.CreateRating(context.Background(), uuid.New(), uuid.New(), &ratingrequest.CreateRatingRequest{Score: 5})

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

func TestGetAverageRatingReturnsSuccessWithEmptyData(t *testing.T) {
	t.Parallel()

	s := newRatingService(t, true)
	res := s.GetAverageRating(context.Background(), uuid.New())

	if !res.Success {
		t.Fatalf("expected success result")
	}
	if res.HttpStatus != http.StatusOK {
		t.Fatalf("expected status 200, got %d", res.HttpStatus)
	}
	if res.Message != "Average rating retrieved successfully" {
		t.Fatalf("unexpected message: %s", res.Message)
	}

	payload, ok := res.Data.(map[string]any)
	if !ok {
		t.Fatalf("expected map data, got %T", res.Data)
	}
	if _, ok := payload["average"]; !ok {
		t.Fatalf("expected average in payload")
	}
	if _, ok := payload["count"]; !ok {
		t.Fatalf("expected count in payload")
	}
}
