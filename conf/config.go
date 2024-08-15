package conf

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	clientv3 "go.etcd.io/etcd/client/v3"
)

func InitConfig() {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath != "" {
		log.Println("Load config file from: ", configPath)
		viper.SetConfigFile(configPath)
	} else {
		log.Println("Load config file from: ./config.yaml")
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(".")
	}

	// 设置默认值
	viper.SetDefault("server.addr", "0.0.0.0:8000")

	// 加载配置文件
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	if viper.GetBool("etcd.enable") {
		// 初始化 etcd 配置
		InitEtcdConfig()
	} else {
		// 监听配置文件变化
		viper.WatchConfig()
		viper.OnConfigChange(func(e fsnotify.Event) {
			fmt.Println("Config file changed:", e.Name)
		})
	}
}

func InitEtcdConfig() {
	// etcd加载配置
	username := viper.GetString("etcd.username")
	password := viper.GetString("etcd.password")
	endpoints := viper.GetStringSlice("etcd.endpoints")
	key := viper.GetString("etcd.key")

	// 创建etcd客户端配置
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints, // etcd服务的地址
		DialTimeout: 5 * time.Second,
		Username:    username,
		Password:    password,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()

	// 读取etcd key配置
	resp, err := cli.Get(context.Background(), key)
	if err != nil {
		log.Fatal(err)
	}

	// 设置viper
	for _, ev := range resp.Kvs {
		if err := viper.MergeConfig(strings.NewReader(string(ev.Value))); err != nil {
			panic(errors.WithMessage(err, "viper read etcd config error"))
		}
	}

	// 监听配置变化
	go func() {
		rch := cli.Watch(context.Background(), key)
		for wresp := range rch {
			for _, ev := range wresp.Events {
				viper.MergeConfig(strings.NewReader(string(ev.Kv.Value)))
			}
		}
	}()
}
