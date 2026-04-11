//go:build integration

package ratingservice

import (
	"context"
	"os"
	"reflect"
	"testing"

	"manga-go/internal/pkg/common"
	"manga-go/internal/pkg/logger"
	comicrepo "manga-go/internal/pkg/repo/comic"
	ratingrepo "manga-go/internal/pkg/repo/rating"
	ratingrequest "manga-go/internal/pkg/request/rating"

	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func newRatingServiceIntegration(t *testing.T) (*RatingService, *gorm.DB) {
	t.Helper()

	dsn := os.Getenv("INTEGRATION_TEST_DATABASE_DSN")
	if dsn == "" {
		t.Skip("INTEGRATION_TEST_DATABASE_DSN is not set")
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect to postgres: %v", err)
	}

	tx := db.Begin()
	if tx.Error != nil {
		t.Fatalf("failed to begin transaction: %v", tx.Error)
	}
	t.Cleanup(func() {
		_ = tx.Rollback().Error
	})

	if err := tx.Exec(`CREATE TABLE ratings (
		id uuid PRIMARY KEY,
		user_id uuid,
		comic_id uuid,
		score INTEGER,
		created_at TIMESTAMPTZ,
		updated_at TIMESTAMPTZ,
		deleted_at TIMESTAMPTZ
	)`).Error; err != nil {
		t.Fatalf("failed to setup schema: %v", err)
	}

	s := &RatingService{
		logger:     logger.NewLogger(),
		ratingRepo: ratingrepo.NewRatingRepo(tx),
		comicRepo:  comicrepo.NewComicRepo(tx),
	}

	return s, tx
}

func ratingPaginationTotalFromData(data any) int64 {
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

func TestRatingServiceIntegrationFullFlow(t *testing.T) {
	s, db := newRatingServiceIntegration(t)
	ctx := context.Background()

	userID := uuid.New()
	comicID := uuid.New()

	createRes := s.CreateRating(ctx, userID, comicID, &ratingrequest.CreateRatingRequest{Score: 3})
	if !createRes.Success {
		t.Fatalf("expected create success, got: %s", createRes.Message)
	}

	upsertRes := s.CreateRating(ctx, userID, comicID, &ratingrequest.CreateRatingRequest{Score: 5})
	if !upsertRes.Success {
		t.Fatalf("expected upsert success, got: %s", upsertRes.Message)
	}

	listRes := s.ListRatings(ctx, userID, &common.Paging{Page: 1, Limit: 10})
	if !listRes.Success {
		t.Fatalf("expected list success, got: %s", listRes.Message)
	}
	if total := ratingPaginationTotalFromData(listRes.Data); total != 1 {
		t.Fatalf("expected total 1, got %d", total)
	}

	avgRes := s.GetAverageRating(ctx, comicID)
	if !avgRes.Success {
		t.Fatalf("expected average success, got: %s", avgRes.Message)
	}
	avgData, ok := avgRes.Data.(map[string]any)
	if !ok {
		t.Fatalf("expected map data, got %T", avgRes.Data)
	}
	if avgData["count"].(int64) != 1 {
		t.Fatalf("expected count 1, got %v", avgData["count"])
	}
	if avgData["average"].(float64) != 5 {
		t.Fatalf("expected average 5, got %v", avgData["average"])
	}

	var ratingID uuid.UUID
	if err := db.Raw("SELECT id FROM ratings WHERE user_id = ? AND comic_id = ? AND deleted_at IS NULL", userID, comicID).Scan(&ratingID).Error; err != nil {
		t.Fatalf("failed to query rating id: %v", err)
	}
	if ratingID == uuid.Nil {
		t.Fatalf("expected persisted rating id")
	}

	updateRes := s.UpdateRating(ctx, userID, ratingID, &ratingrequest.UpdateRatingRequest{Score: 4})
	if !updateRes.Success {
		t.Fatalf("expected update success, got: %s", updateRes.Message)
	}

	deleteRes := s.DeleteRating(ctx, userID, ratingID)
	if !deleteRes.Success {
		t.Fatalf("expected delete success, got: %s", deleteRes.Message)
	}

	notFoundRes := s.UpdateRating(ctx, userID, ratingID, &ratingrequest.UpdateRatingRequest{Score: 2})
	if notFoundRes.Success {
		t.Fatalf("expected not found after soft delete")
	}
	if notFoundRes.Message != "Rating not found" {
		t.Fatalf("unexpected message: %s", notFoundRes.Message)
	}
}
