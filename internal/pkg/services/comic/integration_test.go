//go:build integration

package comicservice

import (
	"context"
	"os"
	"testing"
	"time"

	"manga-go/internal/pkg/constant"
	"manga-go/internal/pkg/logger"
	comicrepo "manga-go/internal/pkg/repo/comic"
	genrerepo "manga-go/internal/pkg/repo/genre"
	tagrepo "manga-go/internal/pkg/repo/tag"
	usercomicreadrepo "manga-go/internal/pkg/repo/user_comic_read"
	comicrequest "manga-go/internal/pkg/request/comic"

	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func newComicServiceIntegration(t *testing.T) (*ComicService, *gorm.DB) {
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

	ddl := []string{
		`CREATE TABLE comics (
			id uuid PRIMARY KEY,
			title TEXT,
			slug TEXT,
			alternative_titles jsonb,
			description TEXT,
			thumbnail TEXT,
			banner TEXT,
			type TEXT,
			status TEXT,
			age_rating TEXT,
			is_published BOOLEAN,
			is_hot BOOLEAN,
			is_featured BOOLEAN,
			published_year INTEGER,
			last_chapter_at TIMESTAMPTZ,
			artist_id uuid,
			created_at TIMESTAMPTZ,
			updated_at TIMESTAMPTZ,
			deleted_at TIMESTAMPTZ
		)`,
	}

	for _, stmt := range ddl {
		if err := tx.Exec(stmt).Error; err != nil {
			t.Fatalf("failed to setup schema: %v", err)
		}
	}

	s := &ComicService{
		logger:            logger.NewLogger(),
		comicRepo:         comicrepo.NewComicRepo(tx),
		genreRepo:         genrerepo.NewGenreRepo(tx),
		tagRepo:           tagrepo.NewTagRepo(tx, nil),
		userComicReadRepo: usercomicreadrepo.NewUserComicReadRepo(tx),
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
