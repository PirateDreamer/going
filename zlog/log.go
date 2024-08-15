package zlog

import (
	"context"

	"github.com/PirateDreamer/going/comm"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// print、write file、upload

var Logger *zap.Logger

func InitZlog() {
	var err error
	if viper.GetBool("dev") {
		Logger, err = zap.NewDevelopment()
	} else {
		Logger, err = zap.NewProduction()
	}
	if err != nil {
		panic(err)
	}
}

func Log(c context.Context) *zap.Logger {
	return Logger.With(zap.String("req_id", comm.GetReqId(c)))
}
