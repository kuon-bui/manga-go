package translationgrouproute

import (
	"manga-go/internal/app/api/common/response"
	translationgrouprequest "manga-go/internal/pkg/request/translation_group"

	"github.com/gin-gonic/gin"
)

// @Summary      Update translation group
// @Description  Update translation group information
// @Tags         TranslationGroup
// @Accept       json
// @Produce      json
// @Param        slug  path      string                                      true  "Translation group slug"
// @Param        body  body      translationgrouprequest.UpdateTranslationGroupRequest  true  "Translation group update request"
// @Success      200   {object}  response.Response
// @Failure      400   {object}  response.Response
// @Failure      401   {object}  response.Response
// @Failure      404   {object}  response.Response
// @Router       /translation-groups/{slug} [put]
// @Security     AccessToken
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
