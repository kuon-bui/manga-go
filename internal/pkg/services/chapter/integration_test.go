//go:build integration

package chapterserivce

import (
	"context"
	"os"
	"reflect"
	"testing"
	"time"

	"manga-go/internal/pkg/common"
	"manga-go/internal/pkg/logger"
	chapterrepo "manga-go/internal/pkg/repo/chapter"
	comicrepo "manga-go/internal/pkg/repo/comic"
	readingprogressrepo "manga-go/internal/pkg/repo/reading_progress"
	usercomicreadrepo "manga-go/internal/pkg/repo/user_comic_read"
	chapterrequest "manga-go/internal/pkg/request/chapter"

	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func newChapterServiceIntegration(t *testing.T) (*ChapterService, *gorm.DB) {
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
		`CREATE TABLE chapters (
			id uuid PRIMARY KEY,
			comic_id uuid,
			number TEXT,
			chapter_idx INTEGER,
			title TEXT,
			slug TEXT,
			is_published BOOLEAN,
			created_at TIMESTAMPTZ,
			updated_at TIMESTAMPTZ,
			deleted_at TIMESTAMPTZ
		)`,
		`CREATE TABLE pages (
			id uuid PRIMARY KEY,
			chapter_id uuid,
			page_number INTEGER,
			image_url TEXT,
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

	s := &ChapterService{
		logger:              logger.NewLogger(),
		chapterRepo:         chapterrepo.NewChapterRepo(tx, nil),
		comicRepo:           comicrepo.NewComicRepo(tx),
		readingProgressRepo: readingprogressrepo.NewReadingProgressRepo(tx),
		userComicReadRepo:   usercomicreadrepo.NewUserComicReadRepo(tx),
	}

	return s, tx
}

func paginationTotalFromResultData(data any) int64 {
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

func TestChapterServiceIntegrationListGetPublishFlow(t *testing.T) {
	s, db := newChapterServiceIntegration(t)
	ctx := context.Background()
	now := time.Now()

	comicID := uuid.New()
	chapterID := uuid.New()
	pageID := uuid.New()

	if err := db.Exec(
		"INSERT INTO chapters (id, comic_id, number, chapter_idx, title, slug, is_published, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)",
		chapterID,
		comicID,
		"1",
		0,
		"Chapter 1",
		"chapter-1",
		false,
		now,
		now,
	).Error; err != nil {
		t.Fatalf("failed to seed chapter: %v", err)
	}

	if err := db.Exec(
		"INSERT INTO pages (id, chapter_id, page_number, image_url, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)",
		pageID,
		chapterID,
		1,
		"https://example.com/page-1.jpg",
		now,
		now,
	).Error; err != nil {
		t.Fatalf("failed to seed page: %v", err)
	}

	comicCtx := common.SetComicIdToContext(ctx, comicID)

	listRes := s.ListChapters(comicCtx, &common.Paging{Page: 1, Limit: 10})
	if !listRes.Success {
		t.Fatalf("expected list success, got: %s", listRes.Message)
	}
	if total := paginationTotalFromResultData(listRes.Data); total != 1 {
		t.Fatalf("expected total 1, got %d", total)
	}

	getRes := s.GetChapter(comicCtx, "chapter-1")
	if !getRes.Success {
		t.Fatalf("expected get success, got: %s", getRes.Message)
	}

	publishRes := s.PublishChapter(comicCtx, "chapter-1", &chapterrequest.PublishChapterRequest{IsPublished: true})
	if !publishRes.Success {
		t.Fatalf("expected publish success, got: %s", publishRes.Message)
	}

	var isPublished bool
	if err := db.Raw("SELECT is_published FROM chapters WHERE id = ?", chapterID).Scan(&isPublished).Error; err != nil {
		t.Fatalf("failed to query chapter publish status: %v", err)
	}
	if !isPublished {
		t.Fatalf("expected chapter to be published")
	}
}
