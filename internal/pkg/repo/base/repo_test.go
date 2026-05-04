package base

import (
	"strings"
	"testing"

	"manga-go/internal/pkg/common"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type testModel struct{}

func (testModel) TableName() string { return "test_models" }

func newDryRunDB(t *testing.T) *gorm.DB {
	t.Helper()

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{DryRun: true})
	if err != nil {
		t.Fatalf("failed to init gorm dry-run db: %v", err)
	}

	return db
}

func TestToClauseColumns(t *testing.T) {
	cols := toClauseColumns([]string{"id", "name"})

	if len(cols) != 2 {
		t.Fatalf("expected len = 2, got %d", len(cols))
	}
	if cols[0].Name != "id" || cols[0].Table != clause.CurrentTable {
		t.Fatalf("unexpected first column: %#v", cols[0])
	}
	if cols[1].Name != "name" || cols[1].Table != clause.CurrentTable {
		t.Fatalf("unexpected second column: %#v", cols[1])
	}
}

func TestWithPaginateNilPagingDoesNotApplyLimitOffset(t *testing.T) {
	repo := &BaseRepository[testModel]{}
	db := newDryRunDB(t)

	sql := db.ToSQL(func(tx *gorm.DB) *gorm.DB {
		scoped := repo.WithPaginate(nil)(tx.Model(&testModel{}))
		return scoped.Find(&[]testModel{})
	})

	lowerSQL := strings.ToLower(sql)
	if strings.Contains(lowerSQL, " limit ") || strings.Contains(lowerSQL, " offset ") {
		t.Fatalf("expected no limit/offset in sql, got %s", sql)
	}
}

func TestWithPaginateAppliesLimitOffset(t *testing.T) {
	repo := &BaseRepository[testModel]{}
	db := newDryRunDB(t)

	paging := &common.Paging{Page: 2, Limit: 10}
	sql := db.ToSQL(func(tx *gorm.DB) *gorm.DB {
		scoped := repo.WithPaginate(paging)(tx.Model(&testModel{}))
		return scoped.Find(&[]testModel{})
	})

	lowerSQL := strings.ToLower(sql)
	if !strings.Contains(lowerSQL, " limit 10") {
		t.Fatalf("expected sql to include limit 10, got %s", sql)
	}
	if !strings.Contains(lowerSQL, " offset 10") {
		t.Fatalf("expected sql to include offset 10, got %s", sql)
	}
}
