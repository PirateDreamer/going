package going

import (
	"github.com/PirateDreamer/going/conf"
	"github.com/PirateDreamer/going/gormx"
	"github.com/PirateDreamer/going/gredis"
	"github.com/PirateDreamer/going/zlog"
)

func InitService() {
	conf.InitConfig(nil)
	zlog.InitZlog()
	gormx.InitMysql()
	gredis.InitRedis()
	InitHttp()
}
