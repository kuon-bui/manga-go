package translationgrouproute

import (
	"github.com/gin-gonic/gin"
)

func (h *TranslationGroupHandler) deleteTranslationGroup(c *gin.Context) {
	slug := c.Param("slug")

	result := h.translationGroupService.DeleteTranslationGroup(c.Request.Context(), slug)
	result.ResponseResult(c)
}
