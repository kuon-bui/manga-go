package fileroute

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestUploadImageMissingFile(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := &FileHandler{}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req := httptest.NewRequest(http.MethodPost, "/files/upload", nil)
	req.Header.Set("Content-Type", "multipart/form-data")
	c.Request = req

	h.uploadImage(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}

func TestGetPresignURLInvalidFilename(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := &FileHandler{}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/files/presign/", nil)
	c.Params = gin.Params{{Key: "filename", Value: ""}}

	h.getPresignURL(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}

func TestGetFileContentInvalidFilename(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := &FileHandler{}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/files/content/", nil)
	c.Params = gin.Params{{Key: "filename", Value: ""}}

	h.getFileContent(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}
