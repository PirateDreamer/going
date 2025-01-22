package main

import (
	"google.golang.org/grpc"
	"net"
	"server/gen/go/mario/user"
	"server/logic"
)

func main() {
	l, err := net.Listen("tcp", ":8000")
	if err != nil {
		panic(err)
	}

	grpcServer := grpc.NewServer()
	user.RegisterUserServiceServer(grpcServer, logic.NewUserServiceServer())

	if err := grpcServer.Serve(l); err != nil {
		panic(err)
	}

	// 注册地址到sidecar，不断上报地址
}
