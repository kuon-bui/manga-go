package commentroute

import (
	"manga-go/internal/app/api/common/response"
	commentrequest "manga-go/internal/pkg/request/comment"
	"manga-go/internal/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// @Summary      Handle reaction for a comment
// @Description  Add or update a reaction for a specific comment
// @Tags         Comment
// @Accept       json
// @Produce      application/json
// @Param        id  path      string  true  "Comment ID"
// @Param        request body commentrequest.AddReactionRequest true "Reaction data"
// @Success      200  {object}  response.Result
// @Failure      400  {object}  response.Result
// @Failure      401  {object}  response.Result
// @Failure      404  {object}  response.Result
// @Failure      500  {object}  response.Result
// @Security     AccessToken
// @Router       /comments/{id}/reactions [post]
func (h *CommentHandler) handleReaction(c *gin.Context) {
	var req commentrequest.AddReactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ResultInvalidRequestErr(err).ResponseResult(c)
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.ResultError("Invalid id").ResponseResult(c)
		return
	}

	user, _ := utils.GetCurrentUserFromGinContext(c)
	result := h.commentService.HandleReaction(c.Request.Context(), user, id, &req)
	result.ResponseResult(c)
}
