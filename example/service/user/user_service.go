package user

import "github.com/PirateDreamer/going/example/internal/repo"

type UserService struct {
	UserRepo repo.UserRepo
}

func NewUserService(userRepo repo.UserRepo) UserService {
	return UserService{
		UserRepo: userRepo,
	}
}
