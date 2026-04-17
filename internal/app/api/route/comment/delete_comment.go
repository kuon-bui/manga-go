package commentroute

import (
	"manga-go/internal/app/api/common/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// @Summary      Delete comment
// @Description  Delete (soft delete) a comment by its ID
// @Tags         Comment
// @Accept       json
// @Produce      json
// @Param        id  path      string  true  "Comment ID"
// @Success      200  {object}  response.Result
// @Failure      400  {object}  response.Result
// @Failure      401  {object}  response.Result
// @Failure      404  {object}  response.Result
// @Failure      500  {object}  response.Result
// @Security     AccessToken
// @Router       /comments/{id} [delete]
func (h *CommentHandler) deleteComment(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.ResultError("Invalid id").ResponseResult(c)
		return
	}

	result := h.commentService.DeleteComment(c.Request.Context(), id)
	result.ResponseResult(c)
}
