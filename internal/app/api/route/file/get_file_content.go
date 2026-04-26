package fileroute

import (
	"mime"
	"net/http"
	"path/filepath"
	"strings"

	"manga-go/internal/app/api/common/response"
	filerequest "manga-go/internal/pkg/request/file"

	"github.com/gin-gonic/gin"
)

// @Summary      Get file content
// @Description  Download file content by path
// @Tags         File
// @Accept       json
// @Produce      */*
// @Param        filename  path      string  true  "File path"
// @Param        variant   query     string  false  "Variant: economy|small|clear|sharp"
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

	var req filerequest.GetFileContentRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.ResultInvalidRequestErr(err).ResponseResult(c)
		return
	}

	fileContent, resolvedKey, err := h.fileService.GetFileByVariant(c.Request.Context(), filename, req.Variant)
	if err != nil {
		if strings.EqualFold(err.Error(), "invalid filename") || strings.EqualFold(err.Error(), "invalid variant") {
			response.ResultError(err.Error()).ResponseResult(c)
			return
		}

		response.ResultErrInternal(err).ResponseResult(c)
		return
	}

	contentType := mime.TypeByExtension(filepath.Ext(resolvedKey))
	if contentType == "" {
		contentType = http.DetectContentType(fileContent)
	}

	c.Header("Cache-Control", "public, max-age=604800")
	c.Data(http.StatusOK, contentType, fileContent)
}
