package zlog

import (
	"context"
	"os"

	"github.com/PirateDreamer/going/comm"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// print、write file、upload

var Logger *zap.Logger

func InitZlog() {
	var err error
	if viper.GetBool("dev") {
		Logger, err = zap.NewDevelopment()
	} else {
		if viper.GetString("log.output") == "file" {
			if viper.GetString("log.file") == "" {
				viper.Set("log.file", "./log_file.log")
			}

			encoderConfig := zap.NewProductionEncoderConfig()
			encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
			core := zapcore.NewCore(
				zapcore.NewJSONEncoder(encoderConfig),
				zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(&lumberjack.Logger{
					Filename:   viper.GetString("log.file"), // 日志文件路径
					MaxSize:    10,                          // 每个日志文件的最大大小（以MB为单位）
					MaxBackups: 5,                           // 保留的旧日志文件的最大数量
					MaxAge:     30,                          // 最多保留的天数
					Compress:   true,                        // 是否压缩旧日志文件
				})),
				zap.NewAtomicLevelAt(zap.InfoLevel),
			)

			Logger = zap.New(core)
		} else {
			Logger, err = zap.NewProduction()
		}
	}
	if err != nil {
		panic(err)
	}
}

func Log(c context.Context) *zap.Logger {
	return Logger.With(zap.String("req_id", comm.GetReqId(c)))
}
