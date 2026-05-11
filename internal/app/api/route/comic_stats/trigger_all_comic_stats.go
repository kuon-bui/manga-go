package comicstatsroute

import (
	"manga-go/internal/app/api/common/response"

	"github.com/gin-gonic/gin"
)

// @Summary      Trigger comic stats update for all comics
// @Description  Trigger recomputation of stats for all comics in the system
// @Tags         Comic Stats
// @Accept       json
// @Produce      json
// @Success      200  {object}  response.Response
// @Failure      401  {object}  response.Response
// @Failure      500  {object}  response.Response
// @Security     AccessToken
// @Router       /admin/comic-stats/trigger-all [post]
func (h *ComicStatsHandler) triggerAllComicStats(c *gin.Context) {
	count, err := h.comicStatsService.TriggerUpdateForAll(c.Request.Context())
	if err != nil {
		response.ResultErrInternal(err).ResponseResult(c)
		return
	}

	response.ResultSuccess("Comic stats update tasks enqueued", map[string]int{"totalComics": count}).ResponseResult(c)
}
