package userroute

import (
	"manga-go/internal/pkg/utils"

	"github.com/gin-gonic/gin"
)

func (h *userHandler) me(c *gin.Context) {
	user, _ := utils.GetCurrentUserFromGinContext(c)
	c.JSON(200, gin.H{
		"user": user,
	})
}
