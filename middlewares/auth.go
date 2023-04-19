package middlewares

import (
	"Im/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

// 检查用户是否登录
func AuthCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("token")
		userClaims, err := utils.AnalyToken(token)
		if err != nil {
			c.Abort()
			c.JSON(http.StatusOK, gin.H{
				"code":    -1,
				"message": "用户认证失败",
			})
			return
		}
		c.Set("user_claims", userClaims)
		c.Next()
	}
}
