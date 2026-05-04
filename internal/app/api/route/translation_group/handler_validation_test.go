package translationgrouproute

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func newTranslationGroupCtx(method, path string, body string) (*gin.Context, *httptest.ResponseRecorder) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req := httptest.NewRequest(method, path, bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	c.Request = req
	return c, w
}

func TestCreateTranslationGroupInvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := &TranslationGroupHandler{}
	c, w := newTranslationGroupCtx(http.MethodPost, "/translation-groups", "{")

	h.createTranslationGroup(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}
