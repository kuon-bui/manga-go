//go:build integration

package comicservice

import (
	"context"
	"testing"
	"time"

	"manga-go/internal/pkg/constant"
	"manga-go/internal/pkg/logger"
	comicrepo "manga-go/internal/pkg/repo/comic"
	genrerepo "manga-go/internal/pkg/repo/genre"
	tagrepo "manga-go/internal/pkg/repo/tag"
	usercomicreadrepo "manga-go/internal/pkg/repo/user_comic_read"
	comicrequest "manga-go/internal/pkg/request/comic"
	"manga-go/internal/pkg/testutil"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func newComicServiceIntegration(t *testing.T) (*ComicService, *gorm.DB) {
	t.Helper()

	db := testutil.NewSQLiteDB(t)
	tx := db.Begin()
	if tx.Error != nil {
		t.Fatalf("failed to begin transaction: %v", tx.Error)
	}
	t.Cleanup(func() {
		_ = tx.Rollback().Error
	})

	testutil.MustSyncSchemas(t, tx,
		&testutil.Comic{},
		&testutil.Chapter{},
		&testutil.Rating{},
		&testutil.ComicFollow{},
		&testutil.UserComicRead{},
	)

	s := &ComicService{
		logger:            logger.NewLogger(),
		comicRepo:         comicrepo.NewComicRepo(tx),
		genreRepo:         genrerepo.NewGenreRepo(tx),
		tagRepo:           tagrepo.NewTagRepo(tx, nil),
		userComicReadRepo: usercomicreadrepo.NewUserComicReadRepo(tx),
		gormDb:            tx,
	}

	return s, tx
}

func TestComicServiceIntegrationStatusPublishDeleteFlow(t *testing.T) {
	s, db := newComicServiceIntegration(t)
	ctx := context.Background()
	now := time.Now()

	comicID := uuid.New()
	comicSlug := "integration-comic"

	if err := db.Exec(
		"INSERT INTO comics (id, title, slug, type, status, age_rating, is_published, is_hot, is_featured, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		comicID,
		"Integration Comic",
		comicSlug,
		constant.ComicTypeManga,
		constant.ComicStatusOngoing,
		constant.AgeRatingAll,
		false,
		false,
		false,
		now,
		now,
	).Error; err != nil {
		t.Fatalf("failed to seed comic: %v", err)
	}

	statusRes := s.UpdateComicStatus(ctx, comicSlug, &comicrequest.UpdateComicStatusRequest{Status: constant.ComicStatusCompleted})
	if !statusRes.Success {
		t.Fatalf("expected update status success, got: %s", statusRes.Message)
	}

	var status string
	if err := db.Raw("SELECT status FROM comics WHERE id = ?", comicID).Scan(&status).Error; err != nil {
		t.Fatalf("failed to query comic status: %v", err)
	}
	if status != string(constant.ComicStatusCompleted) {
		t.Fatalf("expected status %s, got %s", constant.ComicStatusCompleted, status)
	}

	publishRes := s.PublishComic(ctx, comicSlug, &comicrequest.PublishComicRequest{IsPublished: true})
	if !publishRes.Success {
		t.Fatalf("expected publish success, got: %s", publishRes.Message)
	}

	var isPublished bool
	if err := db.Raw("SELECT is_published FROM comics WHERE id = ?", comicID).Scan(&isPublished).Error; err != nil {
		t.Fatalf("failed to query publish state: %v", err)
	}
	if !isPublished {
		t.Fatalf("expected comic to be published")
	}

	deleteRes := s.DeleteComic(ctx, comicSlug)
	if !deleteRes.Success {
		t.Fatalf("expected delete success, got: %s", deleteRes.Message)
	}

	var deletedAt *time.Time
	if err := db.Raw("SELECT deleted_at FROM comics WHERE id = ?", comicID).Scan(&deletedAt).Error; err != nil {
		t.Fatalf("failed to query deleted_at: %v", err)
	}
	if deletedAt == nil {
		t.Fatalf("expected deleted_at to be set")
	}

	notFoundRes := s.UpdateComicStatus(ctx, comicSlug, &comicrequest.UpdateComicStatusRequest{Status: constant.ComicStatusHiatus})
	if notFoundRes.Success {
		t.Fatalf("expected not found after soft delete")
	}
	if notFoundRes.Message != "Comic not found" {
		t.Fatalf("unexpected message after soft delete: %s", notFoundRes.Message)
	}
}
