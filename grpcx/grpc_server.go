package grpcx

import (
	"log"
	"net"

	"google.golang.org/grpc"
)

type GRPCServer struct {
	Server *grpc.Server
}

func NewGrpcServer(opts ...grpc.ServerOption) *GRPCServer {
	if len(opts) == 0 {
		// 默认配置
		opts = []grpc.ServerOption{
			// grpc.UnaryInterceptor(interceptor.LoggingInterceptor), // 日志拦截器
			grpc.UnaryInterceptor(validationInterceptor), // 参数校验拦截器
		}
	}

	s := grpc.NewServer(opts...)
	return &GRPCServer{Server: s}
}

func (s *GRPCServer) StartGrpcServer(addr string) {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	log.Printf("Starting gRPC server on %s", addr)
	if err := s.Server.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

// Stop 关闭服务
func (s *GRPCServer) StopGrpcServer() {
	s.Server.GracefulStop()
}
