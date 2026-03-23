package genreroute

import (
	"manga-go/internal/app/api/common/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h *GenreHandler) getGenre(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.ResultError("invalid id").ResponseResult(c)
		return
	}

	result := h.genreService.GetGenre(c.Request.Context(), id)
	result.ResponseResult(c)
}
