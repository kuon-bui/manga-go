package fileroute

import (
	"strings"

	"manga-go/internal/app/api/common/response"

	"github.com/gin-gonic/gin"
)

func (h *FileHandler) getPresignURL(c *gin.Context) {
	filename := strings.TrimPrefix(c.Param("filename"), "/")
	if filename == "" {
		response.ResultError("Invalid filename").ResponseResult(c)
		return
	}

	url, err := h.fileService.GeneratePresignedURL(c.Request.Context(), filename)
	if err != nil {
		response.ResultErrInternal(err).ResponseResult(c)
		return
	}

	response.ResultSuccess("Get presigned URL successfully", map[string]any{
		"url":      url,
		"filename": filename,
	}).ResponseResult(c)
}
