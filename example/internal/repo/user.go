package repo

import (
	"context"
	"demo/domain"

	"github.com/PirateDreamer/going/gormx"
)

type UserRepo struct {
}

func NewUserRepo() UserRepo {
	return UserRepo{}
}

// GetUserByAccount 根据手机号或邮箱查询用户信息
func (u UserRepo) GetUserByAccount(c context.Context, account string) (user domain.User, err error) {
	err = gormx.Mysql.WithContext(c).Where("phone = ? or email = ?", account, account).First(&user).Error
	return
}
