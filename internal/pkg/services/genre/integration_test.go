//go:build integration

package genreservice

import (
	"context"
	"reflect"
	"testing"

	"manga-go/internal/pkg/common"
	"manga-go/internal/pkg/logger"
	genrerepo "manga-go/internal/pkg/repo/genre"
	genrerequest "manga-go/internal/pkg/request/genre"
	"manga-go/internal/pkg/testutil"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func newGenreServiceIntegration(t *testing.T) (*GenreService, *gorm.DB) {
	t.Helper()

	db := testutil.NewSQLiteDB(t)
	tx := db.Begin()
	if tx.Error != nil {
		t.Fatalf("failed to begin transaction: %v", tx.Error)
	}
	t.Cleanup(func() {
		_ = tx.Rollback().Error
	})

	testutil.MustSyncSchemas(t, tx, &testutil.Genre{})

	s := &GenreService{
		logger:    logger.NewLogger(),
		genreRepo: genrerepo.NewGenreRepo(tx),
	}

	return s, tx
}

func genrePaginationTotalFromData(data any) int64 {
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

func TestGenreServiceIntegrationFullFlow(t *testing.T) {
	s, db := newGenreServiceIntegration(t)
	ctx := context.Background()

	createRes := s.CreateGenre(ctx, &genrerequest.CreateGenreRequest{
		Name:        "Action",
		Slug:        "action",
		Description: "Action genre",
		Thumbnail:   "thumb.jpg",
	})
	if !createRes.Success {
		t.Fatalf("expected create success, got: %s", createRes.Message)
	}

	listRes := s.ListGenres(ctx, &common.Paging{Page: 1, Limit: 10})
	if !listRes.Success {
		t.Fatalf("expected list success, got: %s", listRes.Message)
	}
	if total := genrePaginationTotalFromData(listRes.Data); total != 1 {
		t.Fatalf("expected total 1, got %d", total)
	}

	getRes := s.GetGenre(ctx, "action")
	if !getRes.Success {
		t.Fatalf("expected get success, got: %s", getRes.Message)
	}

	updateRes := s.UpdateGenre(ctx, "action", &genrerequest.UpdateGenreRequest{
		Name:        "Action+",
		Slug:        "action-plus",
		Description: "Updated",
		Thumbnail:   "thumb-2.jpg",
	})
	if !updateRes.Success {
		t.Fatalf("expected update success, got: %s", updateRes.Message)
	}

	deleteRes := s.DeleteGenre(ctx, "action-plus")
	if !deleteRes.Success {
		t.Fatalf("expected delete success, got: %s", deleteRes.Message)
	}

	notFoundRes := s.GetGenre(ctx, "action-plus")
	if notFoundRes.Success {
		t.Fatalf("expected not found after soft delete")
	}
	if notFoundRes.Message != "Genre not found" {
		t.Fatalf("unexpected message: %s", notFoundRes.Message)
	}

	genreID := testutil.MustReadUUID(t, db, "SELECT id FROM genres WHERE slug = ?", "action-plus")
	if genreID == uuid.Nil {
		t.Fatalf("expected persisted genre id")
	}
}
