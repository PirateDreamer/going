package conf

import (
	"log"
	"testing"
	"time"
)

func TestLocalConfig(t *testing.T) {
	type Mysql struct {
		DbName   string `json:"db_name"`
		Host     string `json:"host"`
		Password string `json:"password"`
		Port     int    `json:"port"`
		User     string `json:"user"`
	}

	type ServerConfig struct {
		Mysql Mysql `json:"mysql"`
	}
	// 加载配置文件
	serverConfig := NewConfig[ServerConfig](func() string {
		return "config_one.yaml"
	}, func() string {
		return "config_two.yaml"
	})
	serverConfig.Load()
	serverConfig.Watch()
	for {
		log.Printf("Mysql: %v", serverConfig.Get().Mysql)
		time.Sleep(time.Second)
	}
}
