//go:build integration

package authorservice

import (
	"context"
	"reflect"
	"testing"

	"manga-go/internal/pkg/common"
	"manga-go/internal/pkg/logger"
	authorrepo "manga-go/internal/pkg/repo/author"
	authorrequest "manga-go/internal/pkg/request/author"
	"manga-go/internal/pkg/testutil"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func newAuthorServiceIntegration(t *testing.T) (*AuthorService, *gorm.DB) {
	t.Helper()

	db := testutil.NewSQLiteDB(t)
	tx := db.Begin()

	if tx.Error != nil {
		t.Fatalf("failed to begin transaction: %v", tx.Error)
	}
	t.Cleanup(func() {
		_ = tx.Rollback().Error
	})

	testutil.MustSyncSchemas(t, tx, &testutil.Author{})

	s := &AuthorService{
		logger:     logger.NewLogger(),
		authorRepo: authorrepo.NewAuthorRepo(tx),
	}

	return s, tx
}

func authorPaginationTotalFromData(data any) int64 {
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

func TestAuthorServiceIntegrationFullFlow(t *testing.T) {
	s, db := newAuthorServiceIntegration(t)
	ctx := context.Background()

	createRes := s.CreateAuthor(ctx, &authorrequest.CreateAuthorRequest{Name: "Eiichiro Oda"})
	if !createRes.Success {
		t.Fatalf("expected create success, got: %s", createRes.Message)
	}

	listRes := s.ListAuthors(ctx, &common.Paging{Page: 1, Limit: 10})
	if !listRes.Success {
		t.Fatalf("expected list success, got: %s", listRes.Message)
	}
	if total := authorPaginationTotalFromData(listRes.Data); total != 1 {
		t.Fatalf("expected total 1, got %d", total)
	}

	authorID := testutil.MustReadUUID(t, db, "SELECT id FROM authors WHERE name = ? AND deleted_at IS NULL", "Eiichiro Oda")
	if authorID == uuid.Nil {
		t.Fatalf("expected persisted author id")
	}

	getRes := s.GetAuthor(ctx, authorID)
	if !getRes.Success {
		t.Fatalf("expected get success, got: %s", getRes.Message)
	}

	updateRes := s.UpdateAuthor(ctx, authorID, &authorrequest.UpdateAuthorRequest{Name: "Oda"})
	if !updateRes.Success {
		t.Fatalf("expected update success, got: %s", updateRes.Message)
	}

	deleteRes := s.DeleteAuthor(ctx, authorID)
	if !deleteRes.Success {
		t.Fatalf("expected delete success, got: %s", deleteRes.Message)
	}

	notFoundRes := s.GetAuthor(ctx, authorID)
	if notFoundRes.Success {
		t.Fatalf("expected not found after soft delete")
	}
	if notFoundRes.Message != "Author not found" {
		t.Fatalf("unexpected message: %s", notFoundRes.Message)
	}
}
