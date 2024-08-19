package going

import (
	"github.com/PirateDreamer/going/conf"
	"github.com/PirateDreamer/going/ginx"
	"github.com/PirateDreamer/going/gormx"
	"github.com/PirateDreamer/going/gredis"
	"github.com/PirateDreamer/going/zlog"
	"github.com/gin-gonic/gin"
)

func InitService() (router *gin.Engine) {
	conf.InitConfig(nil)
	zlog.InitZlog()
	gormx.InitMysql()
	gredis.InitRedis()
	router = ginx.InitHttp()
	return
}
