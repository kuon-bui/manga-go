package commentroute

import (
	"manga-go/internal/app/api/common/response"
	commentrequest "manga-go/internal/pkg/request/comment"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// @Summary      Update comment
// @Description  Update a comment's content
// @Tags         Comment
// @Accept       json
// @Produce      json
// @Param        id    path      string                              true  "Comment ID"
// @Param        body  body      commentrequest.UpdateCommentRequest  true  "Comment update request"
// @Success      200   {object}  response.Result
// @Failure      400   {object}  response.Result
// @Failure      401   {object}  response.Result
// @Failure      404   {object}  response.Result
// @Failure      500   {object}  response.Result
// @Security     AccessToken
// @Router       /comments/{id} [put]
func (h *CommentHandler) updateComment(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.ResultError("Invalid id").ResponseResult(c)
		return
	}

	var req commentrequest.UpdateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ResultInvalidRequestErr(err).ResponseResult(c)
		return
	}

	result := h.commentService.UpdateComment(c.Request.Context(), id, &req)
	result.ResponseResult(c)
}
