package grpcx

import "google.golang.org/grpc"

var GrpcServer *grpc.Server

func NewServer() {
	GrpcServer = grpc.NewServer()
}

// 支持拦截器
func RegisterGrpcServer(f func(*grpc.Server, any), srv any) {
	f(GrpcServer, srv)
}

// 创建服务

// 注册grpc服务

// 启动服务
