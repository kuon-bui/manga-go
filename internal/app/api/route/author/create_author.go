package authorroute

import (
	"manga-go/internal/app/api/common/response"
	authorrequest "manga-go/internal/pkg/request/author"

	"github.com/gin-gonic/gin"
)

func (h *AuthorHandler) createAuthor(c *gin.Context) {
	var req authorrequest.CreateAuthorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ResultInvalidRequestErr(err).ResponseResult(c)
		return
	}

	result := h.authorService.CreateAuthor(c.Request.Context(), &req)
	result.ResponseResult(c)
}
