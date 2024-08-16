package user

import "demo/internal/repo"

type UserService struct {
	UserRepo repo.UserRepo
}

func NewUserService(userRepo repo.UserRepo) UserService {
	return UserService{
		UserRepo: userRepo,
	}
}
