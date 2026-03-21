package authorroute

import (
	"github.com/gin-gonic/gin"
)

func (h *AuthorHandler) getAllAuthors(c *gin.Context) {
	result := h.authorService.ListAllAuthors(c.Request.Context())
	result.ResponseResult(c)
}
