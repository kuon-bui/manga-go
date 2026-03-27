package translationgrouproute

import (
	"github.com/gin-gonic/gin"
)

func (h *TranslationGroupHandler) getTranslationGroup(c *gin.Context) {
	slug := c.Param("slug")

	result := h.translationGroupService.GetTranslationGroup(c.Request.Context(), slug)
	result.ResponseResult(c)
}
