package going

import (
	"github.com/PirateDreamer/going/conf"
	"github.com/PirateDreamer/going/ginx"
	"github.com/PirateDreamer/going/gormx"
	"github.com/PirateDreamer/going/gredis"
	"github.com/gin-gonic/gin"
)

func InitService() (router *gin.Engine) {
	conf.InitConfig(nil)
	gormx.InitMysql()
	gredis.InitRedis()
	router = ginx.InitHttp()
	return
}
