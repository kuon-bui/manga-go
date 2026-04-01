package translationgrouproute

import (
	"manga-go/internal/app/api/common/response"
	translationgrouprequest "manga-go/internal/pkg/request/translation_group"
	"manga-go/internal/pkg/utils"

	"github.com/gin-gonic/gin"
)

// @Summary      Create translation group
// @Description  Create a new translation group
// @Tags         TranslationGroup
// @Accept       json
// @Produce      json
// @Param        body  body      translationgrouprequest.CreateTranslationGroupRequest  true  "Translation group creation request"
// @Success      200   {object}  response.Response
// @Failure      400   {object}  response.Response
// @Failure      401   {object}  response.Response
// @Router       /translation-groups [post]
// @Security     AccessToken
func (h *TranslationGroupHandler) createTranslationGroup(c *gin.Context) {
	var req translationgrouprequest.CreateTranslationGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ResultInvalidRequestErr(err).ResponseResult(c)
		return
	}

	currentUser, err := utils.GetCurrentUserFromGinContext(c)
	if err != nil {
		response.ResultUnauthorized().ResponseResult(c)
		return
	}

	result := h.translationGroupService.CreateTranslationGroup(c.Request.Context(), currentUser.ID, &req)
	result.ResponseResult(c)
}
