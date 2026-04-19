package translationgrouproute

import (
	"manga-go/internal/app/api/common/response"
	"manga-go/internal/pkg/common"

	"github.com/gin-gonic/gin"
)

// @Summary      Upload Translation Group logo
// @Description  Upload logo for Translation Group
// @Tags         Translation Group
// @Accept       multipart/form-data
// @Produce      json
// @Param        translationGroupSlug path string true "Translation Group slug"
// @Param        file formData file true "Image file"
// @Success      200    {object}  response.Result
// @Failure      400    {object}  response.Result
// @Failure      401    {object}  response.Result
// @Failure      500    {object}  response.Result
// @Security     AccessToken
// @Router       /translation-groups/{translationGroupSlug}/logo [put]
func (h *TranslationGroupHandler) updateLogo(c *gin.Context) {
	id, ok := common.GetTranslationGroupIdFromContext(c.Request.Context())
	if !ok {
		response.ResultError("invalid translation group id").ResponseResult(c)
		return
	}

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		response.ResultInvalidRequestErr(err).ResponseResult(c)
		return
	}
	defer file.Close()

	contentType := header.Header.Get("Content-Type")

	result := h.translationGroupService.UploadLogo(c.Request.Context(), id, file, header.Size, contentType, header.Filename)
	result.ResponseResult(c)
}
