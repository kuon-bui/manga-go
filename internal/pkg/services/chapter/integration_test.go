//go:build integration

package chapterserivce

import (
	"context"
	"reflect"
	"testing"
	"time"

	"manga-go/internal/pkg/common"
	"manga-go/internal/pkg/logger"
	chapterrepo "manga-go/internal/pkg/repo/chapter"
	comicrepo "manga-go/internal/pkg/repo/comic"
	pagereactionrepo "manga-go/internal/pkg/repo/page_reaction"
	readingprogressrepo "manga-go/internal/pkg/repo/reading_progress"
	usercomicreadrepo "manga-go/internal/pkg/repo/user_comic_read"
	chapterrequest "manga-go/internal/pkg/request/chapter"
	"manga-go/internal/pkg/testutil"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func newChapterServiceIntegration(t *testing.T) (*ChapterService, *gorm.DB) {
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
		&testutil.Page{},
		&testutil.PageReaction{},
	)

	s := &ChapterService{
		logger:              logger.NewLogger(),
		chapterRepo:         chapterrepo.NewChapterRepo(tx, nil),
		comicRepo:           comicrepo.NewComicRepo(tx),
		pageReactionRepo:    pagereactionrepo.NewPageReactionRepo(tx),
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

func TestChapterServiceIntegrationListGetTogglePublishStateFlow(t *testing.T) {
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
		true,
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

	publishRes := s.PublishChapter(comicCtx, "chapter-1", &chapterrequest.PublishChapterRequest{IsPublished: false})
	if !publishRes.Success {
		t.Fatalf("expected publish success, got: %s", publishRes.Message)
	}

	var isPublished bool
	if err := db.Raw("SELECT is_published FROM chapters WHERE id = ?", chapterID).Scan(&isPublished).Error; err != nil {
		t.Fatalf("failed to query chapter publish status: %v", err)
	}
	if isPublished {
		t.Fatalf("expected chapter to be unpublished")
	}
}
