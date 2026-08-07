package authorizationroute

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	authorizationaudit "manga-go/internal/pkg/repo/authorization_audit"
	authorizationadmin "manga-go/internal/pkg/services/authorization_admin"
	"manga-go/internal/pkg/testutil"

	"github.com/gin-gonic/gin"
)

func TestGetAuditLogsUsesStandardPaginationResponse(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := testutil.NewSQLiteDB(t)
	testutil.MustSyncSchemas(t, db, &testutil.AuthorizationAuditLog{})
	handler := &Handler{authorizationAdmin: authorizationadmin.NewService(authorizationadmin.ServiceParams{
		AuditRepo: authorizationaudit.NewRepo(db),
	})}
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Request = httptest.NewRequest(http.MethodGet, "/authorization/audit-logs", nil)

	handler.getAuditLogs(ctx)

	var body struct {
		Message string                     `json:"message"`
		Data    map[string]json.RawMessage `json:"data"`
	}
	if err := json.Unmarshal(recorder.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if body.Message != "Authorization audit logs retrieved successfully" {
		t.Fatalf("unexpected message: %q", body.Message)
	}
	if len(body.Data) != 2 {
		t.Fatalf("expected only data and total pagination fields, got %s", body.Data)
	}
	var entries []json.RawMessage
	if err := json.Unmarshal(body.Data["data"], &entries); err != nil {
		t.Fatalf("expected audit-log data array: %v", err)
	}
	var total int64
	if err := json.Unmarshal(body.Data["total"], &total); err != nil {
		t.Fatalf("expected audit-log total: %v", err)
	}
	if len(entries) != 0 || total != 0 {
		t.Fatalf("unexpected audit-log page: entries=%d total=%d", len(entries), total)
	}
}

func TestGetAuditLogsRejectsInvertedDateRange(t *testing.T) {
	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Request = httptest.NewRequest(
		http.MethodGet,
		"/authorization/audit-logs?start_at=2026-08-08T00:00:00Z&end_at=2026-08-07T00:00:00Z",
		nil,
	)
	handler := &Handler{}

	handler.getAuditLogs(ctx)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", recorder.Code)
	}
}

func TestGetAuditLogsRejectsInvalidTargetID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Request = httptest.NewRequest(
		http.MethodGet,
		"/authorization/audit-logs?target_id=invalid",
		nil,
	)
	handler := &Handler{}

	handler.getAuditLogs(ctx)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", recorder.Code)
	}
}
