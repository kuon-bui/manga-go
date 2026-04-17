package commentroute

import (
	"manga-go/internal/app/api/common/response"
	commentrequest "manga-go/internal/pkg/request/comment"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// @Summary      List comment replies
// @Description  Retrieve paginated direct replies for a specific comment
// @Tags         Comment
// @Accept       json
// @Produce      json
// @Param        id     path   string  true   "Comment ID"
// @Param        page   query  int     false  "Page number (default: 1)"
// @Param        limit  query  int     false  "Records per page (default: 20)"
// @Success      200    {object}  response.Result
// @Failure      400    {object}  response.Result
// @Failure      401    {object}  response.Result
// @Failure      500    {object}  response.Result
// @Security     AccessToken
// @Router       /comments/{id}/replies [get]
func (h *CommentHandler) getCommentReplies(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.ResultError("Invalid id").ResponseResult(c)
		return
	}

	var req commentrequest.ListCommentRepliesRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.ResultInvalidRequestErr(err).ResponseResult(c)
		return
	}

	result := h.commentService.ListCommentReplies(c.Request.Context(), id, &req)
	result.ResponseResult(c)
}
