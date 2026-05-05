//go:build integration

package tagservice

import (
	"context"
	"reflect"
	"testing"

	"manga-go/internal/pkg/common"
	"manga-go/internal/pkg/logger"
	tagrepo "manga-go/internal/pkg/repo/tag"
	tagrequest "manga-go/internal/pkg/request/tag"
	"manga-go/internal/pkg/testutil"

	"gorm.io/gorm"
)

func newTagServiceIntegration(t *testing.T) (*TagService, *gorm.DB) {
	t.Helper()

	db := testutil.NewSQLiteDB(t)
	tx := db.Begin()
	if tx.Error != nil {
		t.Fatalf("failed to begin transaction: %v", tx.Error)
	}
	t.Cleanup(func() {
		_ = tx.Rollback().Error
	})

	testutil.MustSyncSchemas(t, tx, &testutil.Tag{})

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
