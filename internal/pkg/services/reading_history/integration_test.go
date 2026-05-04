//go:build integration

package readinghistoryservice

import (
	"context"
	"os"
	"reflect"
	"testing"
	"time"

	"manga-go/internal/pkg/common"
	"manga-go/internal/pkg/logger"
	readinghistoryrepo "manga-go/internal/pkg/repo/reading_history"
	readinghistoryrequest "manga-go/internal/pkg/request/reading_history"

	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func newReadingHistoryServiceIntegration(t *testing.T) (*ReadingHistoryService, *gorm.DB) {
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

	if err := tx.Exec(`CREATE TABLE reading_histories (
		id uuid PRIMARY KEY,
		user_id uuid,
		chapter_id uuid,
		comic_id uuid,
		last_read_at TIMESTAMPTZ,
		created_at TIMESTAMPTZ,
		updated_at TIMESTAMPTZ,
		deleted_at TIMESTAMPTZ
	)`).Error; err != nil {
		t.Fatalf("failed to setup schema: %v", err)
	}

	s := &ReadingHistoryService{
		logger:             logger.NewLogger(),
		readingHistoryRepo: readinghistoryrepo.NewReadingHistoryRepo(tx),
	}

	return s, tx
}

func readingHistoryPaginationTotalFromData(data any) int64 {
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

func TestReadingHistoryServiceIntegrationFullFlow(t *testing.T) {
	s, db := newReadingHistoryServiceIntegration(t)
	ctx := context.Background()

	userID := uuid.New()
	chapterID := uuid.New()
	comicID := uuid.New()

	createRes := s.CreateReadingHistory(ctx, userID, &readinghistoryrequest.CreateReadingHistoryRequest{
		ChapterID: chapterID,
		ComicID:   comicID,
	})
	if !createRes.Success {
		t.Fatalf("expected create success, got: %s", createRes.Message)
	}

	listRes := s.ListReadingHistories(ctx, userID, &common.Paging{Page: 1, Limit: 10})
	if !listRes.Success {
		t.Fatalf("expected list success, got: %s", listRes.Message)
	}
	if total := readingHistoryPaginationTotalFromData(listRes.Data); total != 1 {
		t.Fatalf("expected total 1, got %d", total)
	}

	var historyID uuid.UUID
	if err := db.Raw("SELECT id FROM reading_histories WHERE user_id = ? AND comic_id = ? AND deleted_at IS NULL", userID, comicID).Scan(&historyID).Error; err != nil {
		t.Fatalf("failed to query reading history id: %v", err)
	}
	if historyID == uuid.Nil {
		t.Fatalf("expected persisted reading history id")
	}

	getRes := s.GetReadingHistory(ctx, historyID)
	if !getRes.Success {
		t.Fatalf("expected get success, got: %s", getRes.Message)
	}

	customTime := time.Now().Add(-1 * time.Hour)
	updateRes := s.UpdateReadingHistory(ctx, historyID, &readinghistoryrequest.UpdateReadingHistoryRequest{LastReadAt: &customTime})
	if !updateRes.Success {
		t.Fatalf("expected update success, got: %s", updateRes.Message)
	}

	deleteRes := s.DeleteReadingHistory(ctx, historyID)
	if !deleteRes.Success {
		t.Fatalf("expected delete success, got: %s", deleteRes.Message)
	}

	notFoundRes := s.GetReadingHistory(ctx, historyID)
	if notFoundRes.Success {
		t.Fatalf("expected not found after soft delete")
	}
	if notFoundRes.Message != "ReadingHistory not found" {
		t.Fatalf("unexpected message: %s", notFoundRes.Message)
	}
}
