package fileroute

import (
	"mime"
	"net/http"
	"path/filepath"
	"strings"

	"manga-go/internal/app/api/common/response"

	"github.com/gin-gonic/gin"
)

func (h *FileHandler) getFileContent(c *gin.Context) {
	filename := strings.TrimPrefix(c.Param("filename"), "/")
	if filename == "" {
		response.ResultError("Invalid filename").ResponseResult(c)
		return
	}

	fileContent, err := h.fileService.GetFile(c.Request.Context(), filename)
	if err != nil {
		response.ResultErrInternal(err).ResponseResult(c)
		return
	}

	contentType := mime.TypeByExtension(filepath.Ext(filename))
	if contentType == "" {
		contentType = http.DetectContentType(fileContent)
	}

	c.Data(http.StatusOK, contentType, fileContent)
}
