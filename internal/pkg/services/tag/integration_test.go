//go:build integration

package tagservice

import (
	"context"
	"os"
	"reflect"
	"testing"

	"manga-go/internal/pkg/common"
	"manga-go/internal/pkg/logger"
	tagrepo "manga-go/internal/pkg/repo/tag"
	tagrequest "manga-go/internal/pkg/request/tag"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func newTagServiceIntegration(t *testing.T) (*TagService, *gorm.DB) {
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

	if err := tx.Exec(`CREATE TABLE tags (
		id uuid PRIMARY KEY,
		name TEXT,
		slug TEXT,
		created_at TIMESTAMPTZ,
		updated_at TIMESTAMPTZ,
		deleted_at TIMESTAMPTZ
	)`).Error; err != nil {
		t.Fatalf("failed to setup schema: %v", err)
	}

	s := &TagService{
		logger:  logger.NewLogger(),
		tagRepo: tagrepo.NewTagRepo(tx, nil),
	}

	return s, tx
}

func tagPaginationTotalFromData(data any) int64 {
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

func TestTagServiceIntegrationFullFlow(t *testing.T) {
	s, _ := newTagServiceIntegration(t)
	ctx := context.Background()

	createRes := s.CreateTag(ctx, &tagrequest.CreateTagRequest{
		Name: "Shounen",
		Slug: "shounen",
	})
	if !createRes.Success {
		t.Fatalf("expected create success, got: %s", createRes.Message)
	}

	listRes := s.ListTags(ctx, &common.Paging{Page: 1, Limit: 10})
	if !listRes.Success {
		t.Fatalf("expected list success, got: %s", listRes.Message)
	}
	if total := tagPaginationTotalFromData(listRes.Data); total != 1 {
		t.Fatalf("expected total 1, got %d", total)
	}

	deleteRes := s.DeleteTag(ctx, "shounen")
	if !deleteRes.Success {
		t.Fatalf("expected delete success, got: %s", deleteRes.Message)
	}

	notFoundRes := s.DeleteTag(ctx, "shounen")
	if notFoundRes.Success {
		t.Fatalf("expected not found on second delete")
	}
	if notFoundRes.Message != "Tag not found" {
		t.Fatalf("unexpected message: %s", notFoundRes.Message)
	}
}
