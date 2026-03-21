package genreroute

import (
	"manga-go/internal/app/api/common/response"
	genrerequest "manga-go/internal/pkg/request/genre"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h *GenreHandler) updateGenre(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.ResultError("invalid id").ResponseResult(c)
		return
	}

	var req genrerequest.UpdateGenreRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ResultInvalidRequestErr(err).ResponseResult(c)
		return
	}

	result := h.genreService.UpdateGenre(c.Request.Context(), id, &req)
	result.ResponseResult(c)
}
