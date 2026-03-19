package userroute

import "github.com/gin-gonic/gin"

func (h *userHandler) logout(c *gin.Context) {
	c.JSON(200, gin.H{
		"msg": "logout",
	})
}
