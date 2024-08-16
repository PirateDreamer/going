package rest

import (
	"context"

	"github.com/PirateDreamer/going/example/internal/repo"
	"github.com/PirateDreamer/going/example/service/user"
	"github.com/PirateDreamer/going/ginx"
	"github.com/gin-gonic/gin"
)

func InitUserIoc() {
	ginx.Container.Provide(repo.NewUserRepo)
	ginx.Container.Provide(user.NewUserService)
	ginx.Container.Provide(NewUserRest)
}

func InitUserRest() {
	ginx.Container.Invoke(func(userRest UserRest) {
		ginx.R.POST("/user/login", ginx.Run(userRest.Login))
	})
}

type UserRest struct {
	UserService user.UserService
}

func NewUserRest(userService user.UserService) UserRest {
	return UserRest{
		UserService: userService,
	}
}

type LoginReq struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

func (u UserRest) Login(ctx context.Context, c *gin.Context, param LoginReq) (resp *ginx.Empty, err error) {
	// logic
	return
}
