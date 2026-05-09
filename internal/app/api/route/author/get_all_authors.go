package authorroute

import (
	"github.com/gin-gonic/gin"
)

// @Summary      Get all authors
// @Description  Get all authors without pagination
// @Tags         Author
// @Accept       json
// @Produce      json
// @Success      200  {object}  response.Result
// @Failure      401  {object}  response.Result
// @Router       /authors/all [get]
// @Security     AccessToken
func (h *AuthorHandler) getAllAuthors(c *gin.Context) {
	result := h.authorService.ListAllAuthors(c.Request.Context())
	result.ResponseResult(c)
}
