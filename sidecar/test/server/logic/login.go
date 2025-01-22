package logic

import (
	"context"
	"server/gen/go/mario/user"
)

func (s UserServiceServer) Login(ctx context.Context, req *user.LoginRequest) (resp *user.LoginResponse, err error) {
	resp = &user.LoginResponse{
		Token: "token",
	}
	return
}
