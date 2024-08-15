package main

import (
	"context"
	"fmt"

	"github.com/PirateDreamer/going/conf"
	"github.com/PirateDreamer/going/ginx"
	"github.com/PirateDreamer/going/zlog"
	"github.com/gin-gonic/gin"
)

func main() {
	conf.InitConfig()
	zlog.InitZlog()
	// gormx.InitMysql()

	r := gin.Default()

	r.Use(ginx.TraceMiddleware())

	r.GET("/ping", ginx.Run(sayHello))

	r.Run()
}

type SayHelloReq struct {
	Name string `form:"name"`
}

type SayHelloResp struct {
	Result string `json:"result"`
}

func sayHello(ctx context.Context, c *gin.Context, param SayHelloReq) (resp *SayHelloResp, err error) {
	zlog.Log(ctx).Info(param.Name)
	resp = &SayHelloResp{
		Result: fmt.Sprintf("hello %s", param.Name),
	}
	return
}
