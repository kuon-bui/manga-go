package translationgrouproute

import (
	"manga-go/internal/app/api/common/response"
	translationgrouprequest "manga-go/internal/pkg/request/translation_group"

	"github.com/gin-gonic/gin"
)

// @Summary      Transfer translation group ownership
// @Description  Transfer ownership of a translation group to another user
// @Tags         TranslationGroup
// @Accept       json
// @Produce      json
// @Param        slug  path      string                                    true  "Translation group slug"
// @Param        body  body      translationgrouprequest.TransferOwnershipRequest  true  "Transfer ownership request"
// @Success      200   {object}  response.Response
// @Failure      400   {object}  response.Response
// @Failure      401   {object}  response.Response
// @Failure      404   {object}  response.Response
// @Router       /translation-groups/{slug}/transfer-ownership [put]
// @Security     AccessToken
func (h *TranslationGroupHandler) transferOwnership(c *gin.Context) {
	slug := c.Param("translationGroupSlug")

	var req translationgrouprequest.TransferOwnershipRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ResultInvalidRequestErr(err).ResponseResult(c)
		return
	}

	result := h.translationGroupService.TransferOwnership(c.Request.Context(), slug, &req)
	result.ResponseResult(c)
}
