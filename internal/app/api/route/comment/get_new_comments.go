package commentroute

import (
	"manga-go/internal/app/api/common/response"
	commentrequest "manga-go/internal/pkg/request/comment"

	"github.com/gin-gonic/gin"
)

// @Summary      List newest comments
// @Description  Retrieve paginated newest top-level comments across all comics and chapters
// @Tags         Comment
// @Accept       json
// @Produce      json
// @Param        page   query  int  false  "Page number (default: 1)"
// @Param        limit  query  int  false  "Records per page (default: 20)"
// @Success      200    {object}  response.Result
// @Failure      400    {object}  response.Result
// @Failure      401    {object}  response.Result
// @Failure      500    {object}  response.Result
// @Security     AccessToken
// @Router       /comments/new [get]
func (h *CommentHandler) getNewComments(c *gin.Context) {
	var req commentrequest.ListNewCommentsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.ResultInvalidRequestErr(err).ResponseResult(c)
		return
	}

	result := h.commentService.ListNewComments(c.Request.Context(), &req)
	result.ResponseResult(c)
}
