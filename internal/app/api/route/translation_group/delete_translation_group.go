package translationgrouproute

import (
	"manga-go/internal/app/api/common/response"
	"manga-go/internal/pkg/utils"

	"github.com/gin-gonic/gin"
)

func (h *TranslationGroupHandler) deleteTranslationGroup(c *gin.Context) {
	slug := c.Param("slug")

	currentUser, err := utils.GetCurrentUserFromGinContext(c)
	if err != nil {
		response.ResultUnauthorized().ResponseResult(c)
		return
	}

	result := h.translationGroupService.DeleteTranslationGroup(c.Request.Context(), currentUser.ID, slug)
	result.ResponseResult(c)
}
