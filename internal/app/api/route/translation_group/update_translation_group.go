package translationgrouproute

import (
	"manga-go/internal/app/api/common/response"
	translationgrouprequest "manga-go/internal/pkg/request/translation_group"

	"github.com/gin-gonic/gin"
)

func (h *TranslationGroupHandler) updateTranslationGroup(c *gin.Context) {
	slug := c.Param("slug")

	var req translationgrouprequest.UpdateTranslationGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ResultInvalidRequestErr(err).ResponseResult(c)
		return
	}

	result := h.translationGroupService.UpdateTranslationGroup(c.Request.Context(), slug, &req)
	result.ResponseResult(c)
}
