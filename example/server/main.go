package main

import (
	"context"
	"log"
	"net"

	userPb "server/gen/user"

	"google.golang.org/grpc"
)

type UserServer struct {
	userPb.UnimplementedUserServiceServer
}

func (s UserServer) Login(ctx context.Context, req *userPb.LoginRequest) (resp *userPb.LoginResponse, err error) {
	return
}

func main() {
	// 网络监听
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	// 初始化grpc服务
	s := grpc.NewServer()
	// 注册服务逻辑
	userPb.RegisterUserServiceServer(s, &UserServer{})
	// 绑定网络
	err = s.Serve(lis)
	if err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
