package main

import (
	"client/gen/go/mario/user"
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

var sidecarAddr = "localhost:8001"

func main() {
	conn, err := grpc.NewClient(sidecarAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	c := user.NewUserServiceClient(conn)

	resp, err := c.Login(context.TODO(), &user.LoginRequest{})
	if err != nil {
		log.Fatal(err)
	}
	log.Println(resp)
}
