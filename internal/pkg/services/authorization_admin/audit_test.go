package authorizationadmin

import (
	"context"
	"testing"
	"time"

	"manga-go/internal/pkg/common"
	"manga-go/internal/pkg/model"
	authorizationaudit "manga-go/internal/pkg/repo/authorization_audit"
	"manga-go/internal/pkg/testutil"

	"github.com/google/uuid"
)

func TestListAuditLogsFiltersActorActionTargetAndDateNewestFirst(t *testing.T) {
	db := testutil.NewSQLiteDB(t)
	testutil.MustSyncSchemas(t, db, &testutil.AuthorizationAuditLog{})
	repo := authorizationaudit.NewRepo(db)
	service := &Service{auditRepo: repo}
	targetID := uuid.New()
	base := time.Date(2026, time.August, 7, 10, 0, 0, 0, time.UTC)

	entries := []*model.AuthorizationAuditLog{
		newAuditEntry("Mai Tran", "mai@example.com", "user.roles_replaced", "user", targetID, base),
		newAuditEntry("Other", "other@example.com", "role.updated", "role", uuid.New(), base.Add(time.Minute)),
		newAuditEntry("Mai Tran", "mai@example.com", "user.roles_replaced", "user", targetID, base.Add(2*time.Minute)),
	}
	for _, entry := range entries {
		if err := repo.AppendTx(db, entry); err != nil {
			t.Fatal(err)
		}
	}

	start := base.Add(-time.Second)
	end := base.Add(3 * time.Minute)
	page, err := service.ListAuditLogs(context.Background(), ListAuditInput{
		Page: 1, Limit: 20, Actor: "MAI@", Action: "user.roles_replaced",
		TargetType: "user", TargetID: &targetID, StartAt: &start, EndAt: &end,
	})
	if err != nil {
		t.Fatal(err)
	}
	if page.Total != 2 || len(page.Data) != 2 {
		t.Fatalf("unexpected page: %#v", page)
	}
	if !page.Data[0].CreatedAt.After(page.Data[1].CreatedAt) {
		t.Fatalf("expected newest first, got %v then %v", page.Data[0].CreatedAt, page.Data[1].CreatedAt)
	}
}

func TestListAuditLogsPaginatesWithoutChangingTotal(t *testing.T) {
	db := testutil.NewSQLiteDB(t)
	testutil.MustSyncSchemas(t, db, &testutil.AuthorizationAuditLog{})
	repo := authorizationaudit.NewRepo(db)
	service := &Service{auditRepo: repo}
	base := time.Date(2026, time.August, 7, 10, 0, 0, 0, time.UTC)
	for i := 0; i < 3; i++ {
		entry := newAuditEntry("Admin", "admin@example.com", "role.created", "role", uuid.New(), base.Add(time.Duration(i)*time.Minute))
		if err := repo.AppendTx(db, entry); err != nil {
			t.Fatal(err)
		}
	}

	page, err := service.ListAuditLogs(context.Background(), ListAuditInput{Page: 2, Limit: 1})
	if err != nil {
		t.Fatal(err)
	}
	if page.Total != 3 || len(page.Data) != 1 || page.Page != 2 || page.Limit != 1 {
		t.Fatalf("unexpected pagination: %#v", page)
	}
}

func newAuditEntry(
	actorName string,
	actorEmail string,
	action string,
	targetType string,
	targetID uuid.UUID,
	createdAt time.Time,
) *model.AuthorizationAuditLog {
	return &model.AuthorizationAuditLog{
		ID:                 uuid.New(),
		ActorUserID:        nil,
		ActorNameSnapshot:  actorName,
		ActorEmailSnapshot: actorEmail,
		Action:             action,
		TargetType:         targetType,
		TargetID:           targetID,
		TargetNameSnapshot: "target",
		Before:             common.JSONMap{"value": "before"},
		After:              common.JSONMap{"value": "after"},
		CreatedAt:          createdAt,
	}
}
