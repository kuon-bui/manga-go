package authorroute

import (
	"manga-go/internal/app/api/common/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// @Summary      Delete author
// @Description  Delete an author by ID
// @Tags         Author
// @Accept       json
// @Produce      json
// @Param        id  path      string  true  "Author ID"
// @Success      200  {object}  response.Response
// @Failure      400  {object}  response.Response
// @Failure      401  {object}  response.Response
// @Failure      404  {object}  response.Response
// @Router       /authors/{id} [delete]
// @Security     AccessToken
func (h *AuthorHandler) deleteAuthor(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.ResultError("invalid id").ResponseResult(c)
		return
	}

	result := h.authorService.DeleteAuthor(c.Request.Context(), id)
	result.ResponseResult(c)
}
