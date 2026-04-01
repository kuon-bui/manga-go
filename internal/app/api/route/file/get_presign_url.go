package fileroute

import (
	"strings"

	"manga-go/internal/app/api/common/response"

	"github.com/gin-gonic/gin"
)

// @Summary      Get presigned URL
// @Description  Generate a presigned URL for accessing a file
// @Tags         File
// @Accept       json
// @Produce      json
// @Param        filename  path      string  true  "File path"
// @Success      200       {object}  response.Response
// @Failure      400       {object}  response.Response
// @Failure      401       {object}  response.Response
// @Failure      500       {object}  response.Response
// @Router       /files/presign/{filename} [get]
// @Security     AccessToken
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
