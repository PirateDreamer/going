package conf

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type ConfigOption func() string

var defaultPath = "config/config.yaml"

// InitConfigs 初始化多个配置文件并返回合并后的viper实例
func InitConfigs(ops ...ConfigOption) *viper.Viper {
	mergedConfig := viper.New()

	for _, op := range ops {
		config := LoadConfig(op)
		// 将每个配置文件的内容合并到主配置中
		for _, key := range config.AllKeys() {
			mergedConfig.Set(key, config.Get(key))
		}

		// 监听配置文件变化
		config.OnConfigChange(func(e fsnotify.Event) {
			fmt.Printf("配置文件 %s 发生变化\n", e.Name)
			// 重新读取配置文件
			if err := config.ReadInConfig(); err != nil {
				fmt.Printf("重新读取配置文件失败: %v\n", err)
				return
			}
			// 更新合并后的配置
			for _, key := range config.AllKeys() {
				mergedConfig.Set(key, config.Get(key))
			}
		})
		config.WatchConfig()
	}

	if viper.GetBool("etcd.enable") {
		// 初始化 etcd 配置
		InitEtcdConfig(mergedConfig)
	}

	return mergedConfig
}

func InitConfig(op ConfigOption) *viper.Viper {

	v := LoadConfig(op)

	viper.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed:", e.Name)
	})

	return v
}

func LoadConfig(op ConfigOption) *viper.Viper {
	path := defaultPath
	if op != nil {
		path = op()
	}

	v := viper.New()

	v.SetConfigFile(path)

	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}

	return v
}

func WithPath(path string) func() string {
	return func() string {
		return ConfigFullPath(path)
	}
}

func WithEnvPath() func() string {
	return func() string {
		path := os.Getenv("CONFIG_PATH")
		return ConfigFullPath(path)
	}
}

func WithDefaultPath() func() string {
	return func() string {
		return defaultPath
	}
}

func ConfigFullPath(path string) string {
	if filepath.IsAbs(path) {
		return path
	}
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return filepath.Join(dir, path)
}

func InitEtcdConfig(v *viper.Viper) {
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
		if err := v.MergeConfig(strings.NewReader(string(ev.Value))); err != nil {
			panic(errors.WithMessage(err, "viper read etcd config error"))
		}
	}

	// 监听配置变化
	go func() {
		rch := cli.Watch(context.Background(), key)
		for wresp := range rch {
			for _, ev := range wresp.Events {
				v.MergeConfig(strings.NewReader(string(ev.Kv.Value)))
			}
		}
	}()
}
