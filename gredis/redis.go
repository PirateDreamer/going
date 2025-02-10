package gredis

import (
	"context"

	"github.com/PirateDreamer/going/zlog"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

var Redis map[int]*redis.Client

func InitRedis() {
	if viper.GetString("redis.addr") == "" && viper.GetBool("redis.disable") {
		return
	}

	Redis = make(map[int]*redis.Client)

	dbs := viper.GetIntSlice("redis.dbs")

	for _, db := range dbs {
		Redis[db] = redis.NewClient(&redis.Options{
			Addr:     viper.GetString("redis.addr"),
			Password: viper.GetString("redis.pass"),
			Username: viper.GetString("redis.username"),
			DB:       db,
		})
		if err := Redis[db].Ping(context.Background()).Err(); err != nil {
			zlog.LogError(context.Background(), "db %d redis init error: %s", db, err.Error())
		}
		zlog.LogInfo(context.Background(), "db %d redis init success", db)
	}
}
