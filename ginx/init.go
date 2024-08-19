package ginx

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func InitHttp() (router *gin.Engine) {
	InitContainer()
	if !viper.GetBool("server.gin_not_default") {
		router = gin.Default()
	} else {
		router = gin.New()
	}
	router.Use(TraceMiddleware())
	R = router.Group("/api")

	AuthR = R.Group("/auth")
	// ginx.AuthR.Use(ginx.AuthMiddleware())
	return
}
