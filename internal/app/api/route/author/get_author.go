package authorroute

import (
	"manga-go/internal/app/api/common/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h *AuthorHandler) getAuthor(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.ResultError("invalid id").ResponseResult(c)
		return
	}

	result := h.authorService.GetAuthor(c.Request.Context(), id)
	result.ResponseResult(c)
}
