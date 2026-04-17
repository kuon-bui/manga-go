package fileroute

import (
	"mime"
	"net/http"
	"path/filepath"
	"strings"

	"manga-go/internal/app/api/common/response"

	"github.com/gin-gonic/gin"
)

// @Summary      Get file content
// @Description  Download file content by path
// @Tags         File
// @Accept       json
// @Produce      */*
// @Param        filename  path      string  true  "File path"
// @Success      200       {file}    string  "File content"
// @Failure      400       {object}  response.Response
// @Failure      500       {object}  response.Response
// @Router       /files/content/{filename} [get]
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
