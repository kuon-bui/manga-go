package translationgrouproute

import (
	"manga-go/internal/app/api/common/response"
	translationgrouprequest "manga-go/internal/pkg/request/translation_group"
	"manga-go/internal/pkg/utils"

	"github.com/gin-gonic/gin"
)

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
