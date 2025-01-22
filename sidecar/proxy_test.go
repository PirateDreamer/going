package sidecar

import (
	"google.golang.org/grpc"
	"log"
	"net"
	"testing"
)

func TestGrpcProxy(t *testing.T) {
	// 启动 gRPC 代理服务器
	lis, err := net.Listen("tcp", ":8001")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// 创建 gRPC 服务器
	grpcServer := grpc.NewServer()

	// 注册通用代理
	proxy := &grpcProxy{
		backendAddr: "localhost:8000", // 替换为后端 gRPC 服务器地址
	}
	grpcServer.RegisterService(&grpc.ServiceDesc{
		ServiceName: "grpcProxy",
		HandlerType: (*grpcProxy)(nil),
		Methods:     []grpc.MethodDesc{},
		Streams: []grpc.StreamDesc{
			{
				StreamName:    "ProxyStream",
				Handler:       proxy.ProxyStream,
				ServerStreams: true,
				ClientStreams: true,
			},
		},
	}, proxy)

	log.Println("gRPC proxy server is running on :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
