package authorizationrevision

import (
	"context"
	"testing"

	"manga-go/internal/pkg/model"
	"manga-go/internal/pkg/testutil"

	"github.com/google/uuid"
)

func TestCurrentInitializesGlobalAndUserRevisions(t *testing.T) {
	db := testutil.NewSQLiteDB(t)
	testutil.MustSyncSchemas(t, db, &model.AuthorizationCacheRevision{})
	repo := NewRepo(db)
	userID := uuid.New()

	global, user, err := repo.Current(context.Background(), userID)
	if err != nil {
		t.Fatal(err)
	}
	if global != 1 || user != 1 {
		t.Fatalf("expected initial revisions 1/1, got %d/%d", global, user)
	}

	var count int64
	if err := db.Model(&model.AuthorizationCacheRevision{}).Count(&count).Error; err != nil {
		t.Fatal(err)
	}
	if count != 2 {
		t.Fatalf("expected global and user rows, got %d", count)
	}
}

func TestBumpMethodsIncrementDurableRevision(t *testing.T) {
	db := testutil.NewSQLiteDB(t)
	testutil.MustSyncSchemas(t, db, &model.AuthorizationCacheRevision{})
	repo := NewRepo(db)
	userID := uuid.New()

	if _, _, err := repo.Current(context.Background(), userID); err != nil {
		t.Fatal(err)
	}
	global, err := repo.BumpGlobalTx(db)
	if err != nil {
		t.Fatal(err)
	}
	user, err := repo.BumpUserTx(db, userID)
	if err != nil {
		t.Fatal(err)
	}
	if global != 2 || user != 2 {
		t.Fatalf("expected bumped revisions 2/2, got %d/%d", global, user)
	}
}
