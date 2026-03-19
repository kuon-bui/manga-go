package userroute

import "github.com/gin-gonic/gin"

func (h *userHandler) renewToken(c *gin.Context) {
	c.JSON(200, gin.H{
		"msg": "renew token",
	})
}
