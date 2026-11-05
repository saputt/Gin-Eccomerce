package response

import "github.com/gin-gonic/gin"

func ErrorRes(c *gin.Context, status int, message string) {
	c.AbortWithStatusJSON(status, gin.H{
		"message": message,
	})
}
