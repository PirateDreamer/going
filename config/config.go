package config

import (
	"fmt"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

func InitConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	// 设置默认值
	viper.SetDefault("server.addr", "0.0.0.0:8000")

	// 加载配置文件
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	// 如果开启 nacos 配置中心，则初始化 nacos 配置，否则监听配置文件变化（不能同时使用 nacos 和监听本地文件）
	if viper.GetBool("nacos.enable") {
		// 初始化 nacos 配置
		InitNacosConfig()
	} else {
		// 监听配置文件变化
		viper.WatchConfig()
		viper.OnConfigChange(func(e fsnotify.Event) {
			fmt.Println("Config file changed:", e.Name)
		})
	}
}

func InitNacosConfig() {
	host := viper.GetString("nacos.host")
	port := viper.GetInt("nacos.port")
	namespace := viper.GetString("nacos.namespace")
	dataId := viper.GetString("nacos.data_id")
	group := viper.GetString("nacos.group")
	username := viper.GetString("nacos.username")
	password := viper.GetString("nacos.password")

	// 服务端配置（如果存在多个服务端配置，会轮询获取）
	sc := []constant.ServerConfig{
		*constant.NewServerConfig(host, uint64(port), constant.WithContextPath("/nacos")),
	}

	// 客户端配置
	cc := *constant.NewClientConfig(
		constant.WithNamespaceId(namespace),
		constant.WithTimeoutMs(5000),
		constant.WithNotLoadCacheAtStart(true),
		constant.WithLogDir("/tmp/nacos/log"),
		constant.WithCacheDir("/tmp/nacos/cache"),
		constant.WithUsername(username),
		constant.WithPassword(password),
	)

	// 创建客户端
	client, err := clients.NewConfigClient(
		vo.NacosClientParam{
			ClientConfig:  &cc,
			ServerConfigs: sc,
		},
	)

	if err != nil {
		panic(errors.WithMessage(err, "create nacos client error"))
	}

	// 获取 nacos 配置
	content, err := client.GetConfig(vo.ConfigParam{
		DataId: dataId,
		Group:  group,
	})
	if err != nil {
		panic(errors.WithMessage(err, "get nacos config error"))
	}

	// 写入 viper
	viper.SetConfigType("yaml")
	if err := viper.MergeConfig(strings.NewReader(content)); err != nil {
		panic(errors.WithMessage(err, "viper read nacos config error"))
	}

	// 监听配置
	err = client.ListenConfig(vo.ConfigParam{
		DataId: dataId,
		Group:  group,
		OnChange: func(namespace, group, dataId, data string) {
			// 合并入 viper
			viper.MergeConfig(strings.NewReader(data))
		},
	})
	if err != nil {
		panic(errors.WithMessage(err, "listen nacos config error"))
	}
}
