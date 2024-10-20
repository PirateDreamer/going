package grpcx

import (
	"log"
	"net"

	"google.golang.org/grpc"
)

type GRPCServer struct {
	Server *grpc.Server
}

func NewGrpcServer(userInterceptors []grpc.UnaryServerInterceptor) *GRPCServer {

	interceptors := []grpc.UnaryServerInterceptor{
		ValidationInterceptor,    // 参数校验拦截器
		PanicRecoveryInterceptor, // panic 恢复拦截器
	}

	// 将用户自定义的拦截器和系统内置的拦截器合并
	interceptors = append(interceptors, userInterceptors...)

	chainUnaryInterceptor := grpc.ChainUnaryInterceptor(interceptors...)

	opts := []grpc.ServerOption{
		chainUnaryInterceptor,
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
