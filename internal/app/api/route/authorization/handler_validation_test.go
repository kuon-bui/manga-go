package authorizationroute

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

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
