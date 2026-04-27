package pageroute

import (
	"manga-go/internal/app/api/common/response"
	pagerequest "manga-go/internal/pkg/request/page"
	"manga-go/internal/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// @Summary      Handle reaction for a page
// @Description  Add or remove a reaction for a specific page
// @Tags         Page
// @Accept       json
// @Produce      application/json
// @Param        pageId   path      string                       true  "Page ID"
// @Param        request  body      pagerequest.AddReactionRequest  true  "Reaction data"
// @Success      200      {object}  response.Result
// @Failure      400      {object}  response.Result
// @Failure      401      {object}  response.Result
// @Failure      404      {object}  response.Result
// @Failure      500      {object}  response.Result
// @Security     AccessToken
// @Router       /pages/{pageId}/reactions [post]
func (h *PageHandler) handleReaction(c *gin.Context) {
	var req pagerequest.AddReactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ResultInvalidRequestErr(err).ResponseResult(c)
		return
	}

	pageId, err := uuid.Parse(c.Param("pageId"))
	if err != nil {
		response.ResultError("Invalid pageId").ResponseResult(c)
		return
	}

	user, _ := utils.GetCurrentUserFromGinContext(c)
	result := h.pageService.HandleReaction(c.Request.Context(), user, pageId, &req)
	result.ResponseResult(c)
}
