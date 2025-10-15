package middleware_test

import (
	"testing"

	"github.com/PirateDreamer/going/ginc/middleware"
	"github.com/gin-gonic/gin"
)

func TestAesApiDataDecrypt(t *testing.T) {
	router := gin.Default()
	router.Use(middleware.AesApiDataDecrypt("1234567890123456", "1234567890123456", 60000))
	router.POST("/ping", func(c *gin.Context) {
		type Ping struct {
			Name string `json:"name"`
		}
		var ping Ping
		if err := c.ShouldBindJSON(&ping); err != nil {
			c.JSON(400, gin.H{
				"message": err.Error(),
			})
			return
		}
		c.JSON(200, gin.H{
			"message": ping.Name,
		})
	})
	router.Run() // 默认监听 0.0.0.0:8080
}
