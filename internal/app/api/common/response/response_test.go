package response

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestResultConflictCarriesStableCode(t *testing.T) {
	result := ResultConflict("ROLE_IN_USE", "role is assigned")

	if result.HttpStatus != http.StatusConflict || result.Code != "ROLE_IN_USE" {
		t.Fatalf("unexpected conflict result: %#v", result)
	}
}

func TestResponseResultWritesStableCode(t *testing.T) {
	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Request = httptest.NewRequest(http.MethodDelete, "/roles/assigned", nil)

	ResultConflict("ROLE_IN_USE", "role is assigned").ResponseResult(ctx)

	var body Response
	if err := json.Unmarshal(recorder.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if body.Code != "ROLE_IN_USE" {
		t.Fatalf("expected response code ROLE_IN_USE, got %#v", body)
	}
}
