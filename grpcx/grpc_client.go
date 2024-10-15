package grpcx

import (
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func GetServerClient[T any](name string, newServerClient func(grpc.ClientConnInterface) T) (T, error) {
	addr := viper.GetString("discovery." + name)

	// 如果自定义服务地址，则使用 k8s 的 DNS 服务发现
	if addr == "" {
		addr = name + ":8000"
	}

	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		var t T
		return t, err
	}

	return newServerClient(conn), nil
}
