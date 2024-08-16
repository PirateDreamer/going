package rest

import (
	"context"
	"demo/domain"
	"demo/internal/repo"
	"demo/service/user"
	"errors"

	"github.com/PirateDreamer/going/ginx"
	"github.com/PirateDreamer/going/xerr"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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
	UserRepo    repo.UserRepo
}

func NewUserRest(userService user.UserService, userRepo repo.UserRepo) UserRest {
	return UserRest{
		UserService: userService,
	}
}

type LoginReq struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

func (u UserRest) Login(ctx context.Context, c *gin.Context, param LoginReq) (resp *domain.User, err error) {
	// logic
	var userInfo domain.User
	if userInfo, err = u.UserRepo.GetUserByAccount(ctx, param.UserName); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = xerr.NewCommBizErr("用户不存在")
		}
		return
	}
	resp = &userInfo
	return
}
