package logic

import "server/gen/go/mario/user"

type UserServiceServer struct {
	user.UnimplementedUserServiceServer
}

func NewUserServiceServer() UserServiceServer {
	return UserServiceServer{}
}
