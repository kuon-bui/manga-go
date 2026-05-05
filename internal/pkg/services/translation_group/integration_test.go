//go:build integration

package translationgroupservice

import (
	"context"
	"reflect"
	"testing"

	"manga-go/internal/pkg/common"
	"manga-go/internal/pkg/logger"
	"manga-go/internal/pkg/model"
	translationgrouprepo "manga-go/internal/pkg/repo/translation_group"
	userrepo "manga-go/internal/pkg/repo/user"
	translationgrouprequest "manga-go/internal/pkg/request/translation_group"
	"manga-go/internal/pkg/testutil"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func newTranslationGroupServiceIntegration(t *testing.T) (*TranslationGroupService, *gorm.DB, uuid.UUID) {
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
		&testutil.TranslationGroup{},
		&testutil.User{},
	)

	ownerID := uuid.New()
	if err := tx.Exec(`INSERT INTO users (id, name, email, password) VALUES (?, ?, ?, ?)`, ownerID, "owner", "owner@example.com", "secret").Error; err != nil {
		t.Fatalf("failed to seed owner user: %v", err)
	}

	s := &TranslationGroupService{
		logger:               logger.NewLogger(),
		translationGroupRepo: translationgrouprepo.NewTranslationGroupRepo(tx, nil),
		userRepo:             userrepo.NewUserRepository(tx, nil),
	}

	return s, tx, ownerID
}

func translationGroupPaginationTotalFromData(data any) int64 {
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

func TestTranslationGroupServiceIntegrationFullFlow(t *testing.T) {
	s, db, ownerID := newTranslationGroupServiceIntegration(t)
	ctx := context.Background()

	createRes := s.CreateTranslationGroup(ctx, ownerID, &translationgrouprequest.CreateTranslationGroupRequest{
		Name: "Akatsuki Team",
		Slug: "akatsuki-team",
	})
	if !createRes.Success {
		t.Fatalf("expected create success, got: %s", createRes.Message)
	}

	listRes := s.ListTranslationGroups(ctx, &common.Paging{Page: 1, Limit: 10})
	if !listRes.Success {
		t.Fatalf("expected list success, got: %s", listRes.Message)
	}
	if total := translationGroupPaginationTotalFromData(listRes.Data); total != 1 {
		t.Fatalf("expected total 1, got %d", total)
	}

	getRes := s.GetTranslationGroup(ctx, "akatsuki-team")
	if !getRes.Success {
		t.Fatalf("expected get success, got: %s", getRes.Message)
	}

	updateRes := s.UpdateTranslationGroup(ctx, "akatsuki-team", &translationgrouprequest.UpdateTranslationGroupRequest{
		Name: "Akatsuki Team Updated",
		Slug: "akatsuki-team-updated",
	})
	if !updateRes.Success {
		t.Fatalf("expected update success, got: %s", updateRes.Message)
	}

	updatedGroupID := testutil.MustReadUUID(t, db, "SELECT id FROM translation_groups WHERE slug = ? AND deleted_at IS NULL", "akatsuki-team-updated")
	if updatedGroupID == uuid.Nil {
		t.Fatalf("expected persisted translation group id")
	}

	ownerTranslationGroupID := testutil.MustReadUUID(t, db, "SELECT translation_group_id FROM users WHERE id = ?", ownerID)
	if ownerTranslationGroupID != updatedGroupID {
		t.Fatalf("expected owner translation_group_id to be %s, got %s", updatedGroupID, ownerTranslationGroupID)
	}

	deleteRes := s.DeleteTranslationGroup(ctx, "akatsuki-team-updated")
	if !deleteRes.Success {
		t.Fatalf("expected delete success, got: %s", deleteRes.Message)
	}

	notFoundRes := s.GetTranslationGroup(ctx, "akatsuki-team-updated")
	if notFoundRes.Success {
		t.Fatalf("expected not found after soft delete")
	}
	if notFoundRes.Message != "TranslationGroup not found" {
		t.Fatalf("unexpected message: %s", notFoundRes.Message)
	}
}

func TestTranslationGroupServiceIntegrationGetTranslationGroupPreloadsRelations(t *testing.T) {
	s, db, ownerID := newTranslationGroupServiceIntegration(t)
	ctx := context.Background()

	createRes := s.CreateTranslationGroup(ctx, ownerID, &translationgrouprequest.CreateTranslationGroupRequest{
		Name: "Naruto Team",
		Slug: "naruto-team",
	})
	if !createRes.Success {
		t.Fatalf("expected create success, got: %s", createRes.Message)
	}

	groupID := testutil.MustReadUUID(t, db, "SELECT id FROM translation_groups WHERE slug = ? AND deleted_at IS NULL", "naruto-team")
	if groupID == uuid.Nil {
		t.Fatalf("expected persisted translation group id")
	}

	memberID := uuid.New()
	if err := db.Exec(
		`INSERT INTO users (id, name, email, password, translation_group_id) VALUES (?, ?, ?, ?, ?)`,
		memberID,
		"member",
		"member2@example.com",
		"secret",
		groupID,
	).Error; err != nil {
		t.Fatalf("failed to seed member user: %v", err)
	}

	getRes := s.GetTranslationGroup(ctx, "naruto-team")
	if !getRes.Success {
		t.Fatalf("expected get success, got: %s", getRes.Message)
	}

	group, ok := getRes.Data.(*model.TranslationGroup)
	if !ok {
		t.Fatalf("expected *model.TranslationGroup data, got %T", getRes.Data)
	}
	if group.Owner == nil {
		t.Fatalf("expected Owner to be preloaded")
	}
	if group.Owner.ID != ownerID {
		t.Fatalf("expected Owner ID %s, got %s", ownerID, group.Owner.ID)
	}

	hasOwnerInMembers := false
	hasNewMemberInMembers := false
	for _, member := range group.Members {
		if member.ID == ownerID {
			hasOwnerInMembers = true
		}
		if member.ID == memberID {
			hasNewMemberInMembers = true
		}
	}
	if !hasOwnerInMembers {
		t.Fatalf("expected owner to be present in Members preload")
	}
	if !hasNewMemberInMembers {
		t.Fatalf("expected seeded member to be present in Members preload")
	}
}

func TestTranslationGroupServiceIntegrationTransferOwnership(t *testing.T) {
	s, db, ownerID := newTranslationGroupServiceIntegration(t)
	ctx := context.Background()

	createRes := s.CreateTranslationGroup(ctx, ownerID, &translationgrouprequest.CreateTranslationGroupRequest{
		Name: "One Piece Team",
		Slug: "one-piece-team",
	})
	if !createRes.Success {
		t.Fatalf("expected create success, got: %s", createRes.Message)
	}

	groupID := testutil.MustReadUUID(t, db, "SELECT id FROM translation_groups WHERE slug = ? AND deleted_at IS NULL", "one-piece-team")
	if groupID == uuid.Nil {
		t.Fatalf("expected persisted translation group id")
	}

	outsiderID := uuid.New()
	if err := db.Exec(
		`INSERT INTO users (id, name, email, password) VALUES (?, ?, ?, ?)`,
		outsiderID,
		"outsider",
		"outsider@example.com",
		"secret",
	).Error; err != nil {
		t.Fatalf("failed to seed outsider user: %v", err)
	}

	nonMemberRes := s.TransferOwnership(ctx, "one-piece-team", &translationgrouprequest.TransferOwnershipRequest{
		NewOwnerID: outsiderID,
	})
	if nonMemberRes.Success {
		t.Fatalf("expected transfer ownership to fail for non-member")
	}
	if nonMemberRes.Message != "New owner must be a member of the translation group" {
		t.Fatalf("unexpected message: %s", nonMemberRes.Message)
	}

	memberID := uuid.New()
	if err := db.Exec(
		`INSERT INTO users (id, name, email, password, translation_group_id) VALUES (?, ?, ?, ?, ?)`,
		memberID,
		"member",
		"member@example.com",
		"secret",
		groupID,
	).Error; err != nil {
		t.Fatalf("failed to seed member user: %v", err)
	}

	transferRes := s.TransferOwnership(ctx, "one-piece-team", &translationgrouprequest.TransferOwnershipRequest{
		NewOwnerID: memberID,
	})
	if !transferRes.Success {
		t.Fatalf("expected transfer ownership success, got: %s", transferRes.Message)
	}

	newOwnerID := testutil.MustReadUUID(t, db, "SELECT owner_id FROM translation_groups WHERE id = ?", groupID)
	if newOwnerID != memberID {
		t.Fatalf("expected owner_id to be %s, got %s", memberID, newOwnerID)
	}
}
