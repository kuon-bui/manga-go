package commentroute

import (
	"manga-go/internal/app/api/common/response"
	commentrequest "manga-go/internal/pkg/request/comment"
	"manga-go/internal/pkg/utils"

	"github.com/gin-gonic/gin"
)

// @Summary      Create comment
// @Description  Create a new comment on a chapter or reply to an existing comment
// @Tags         Comment
// @Accept       json
// @Produce      json
// @Param        body  body      commentrequest.CreateCommentRequest  true  "Comment creation request"
// @Success      200   {object}  response.Result
// @Failure      400   {object}  response.Result
// @Failure      401   {object}  response.Result
// @Failure      500   {object}  response.Result
// @Security     AccessToken
// @Router       /comments [post]
func (h *CommentHandler) createComment(c *gin.Context) {
	user, err := utils.GetCurrentUserFromGinContext(c)
	if err != nil {
		response.ResultUnauthorized().ResponseResult(c)
		return
	}

	var req commentrequest.CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ResultInvalidRequestErr(err).ResponseResult(c)
		return
	}

	result := h.commentService.CreateComment(c.Request.Context(), user.ID, &req)
	result.ResponseResult(c)
}
