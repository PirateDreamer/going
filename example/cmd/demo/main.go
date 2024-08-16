package main

import (
	"flag"

	"demo/api"

	"github.com/PirateDreamer/going/conf"
	"github.com/PirateDreamer/going/ginx"
	"github.com/PirateDreamer/going/gormx"
	"github.com/PirateDreamer/going/zlog"
	"github.com/gin-gonic/gin"
)

var configFile = flag.String("f", "/Users/mac/Desktop/workspace/myGoWorkspace/me/going/example/etc/config.yml", "the config f")

func main() {
	flag.Parse()
	conf.InitConfig(configFile)
	zlog.InitZlog()
	gormx.InitMysql()

	ginx.InitContainer()
	r := gin.Default()
	r.Use(ginx.TraceMiddleware())
	ginx.R = r.Group("/api")
	ginx.InitContainer()

	api.InitApi()
	// ginx.AuthR = ginc.R.Group("/auth")
	// ginx.Use(ginx.AuthMiddleware())

	r.Run()
}
