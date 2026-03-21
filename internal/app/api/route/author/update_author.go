package authorroute

import (
	"manga-go/internal/app/api/common/response"
	authorrequest "manga-go/internal/pkg/request/author"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h *AuthorHandler) updateAuthor(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.ResultError("invalid id").ResponseResult(c)
		return
	}

	var req authorrequest.UpdateAuthorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ResultInvalidRequestErr(err).ResponseResult(c)
		return
	}

	result := h.authorService.UpdateAuthor(c.Request.Context(), id, &req)
	result.ResponseResult(c)
}
