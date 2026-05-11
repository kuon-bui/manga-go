package comicstatsroute

import (
	"manga-go/internal/app/api/common/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// @Summary      Trigger comic stats update
// @Description  Trigger recomputation of stats for a specific comic
// @Tags         Comic Stats
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Comic ID (UUID)"
// @Success      200  {object}  response.Response
// @Failure      400  {object}  response.Response
// @Failure      401  {object}  response.Response
// @Failure      500  {object}  response.Response
// @Security     AccessToken
// @Router       /admin/comic-stats/trigger/{id} [post]
func (h *ComicStatsHandler) triggerComicStats(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.ResultError("Invalid comic id").ResponseResult(c)
		return
	}

	if err := h.comicStatsService.TriggerUpdateForComic(c.Request.Context(), id); err != nil {
		response.ResultErrInternal(err).ResponseResult(c)
		return
	}

	response.ResultSuccess("Comic stats update task enqueued", nil).ResponseResult(c)
}
