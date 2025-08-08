package conf

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type ConfigOption func() string

var defaultPath = "config/config.yaml"

// 加载配置
func Load(ops ...ConfigOption) *viper.Viper {
	mergedConfig := viper.New()

	for _, op := range ops {
		config := loadConfig(op)
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

	// if viper.GetBool("etcd.enable") {
	// 	// 初始化 etcd 配置
	// 	InitEtcdConfig(mergedConfig)
	// }

	return mergedConfig
}

func loadConfig(op ConfigOption) *viper.Viper {
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
