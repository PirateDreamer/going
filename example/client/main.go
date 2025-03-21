package main

import (
	userPb "client/gen/user"
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.NewClient("0.0.0.0:50052")
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	// 创建客户端
	client := userPb.NewUserServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// 调用服务端逻辑
	resp, err := client.Login(ctx, &userPb.LoginRequest{})
	if err != nil {
		log.Fatalf("login error: %v", err)
	}
	log.Printf("login resp: %s", resp.Token)

}
