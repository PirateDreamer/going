package gredis

import (
	"context"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

var Redis *redis.Client

func InitRedis() {
	if viper.GetString("redis.addr") == "" {
		return
	}
	Redis = redis.NewClient(&redis.Options{
		Addr:     viper.GetString("redis.addr"),
		Password: viper.GetString("redis.pass"),
		DB:       viper.GetInt("redis.db"),
	})

	if err := Redis.Ping(context.Background()).Err(); err != nil {
		panic(err)
	}
}
