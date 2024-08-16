package domain

import "time"

type User struct {
	Id            int64      `gorm:"column:id;primary_key;AUTO_INCREMENT" json:"id"`
	Nickname      string     `gorm:"column:nickname;NOT NULL" json:"nickname"`                      // 昵称
	Phone         string     `gorm:"column:phone;NOT NULL" json:"phone"`                            // 手机号
	Email         string     `gorm:"column:email;NOT NULL" json:"email"`                            // 邮箱
	Password      string     `gorm:"column:password;NOT NULL" json:"password"`                      // 密码
	Avatar        string     `gorm:"column:avatar;NOT NULL" json:"avatar"`                          // 头像
	Status        int        `gorm:"column:status;default:0;NOT NULL" json:"status"`                // 状态
	LastLoginTime *time.Time `gorm:"column:last_login_time" json:"last_login_time"`                 // 最后登录时间
	LastLoginIp   string     `gorm:"column:last_login_ip;NOT NULL" json:"last_login_ip"`            // 最后登录ip
	CreatedAt     time.Time  `gorm:"column:created_at;default:CURRENT_TIMESTAMP" json:"created_at"` // 创建时间
	UpdatedAt     time.Time  `gorm:"column:updated_at;default:CURRENT_TIMESTAMP" json:"updated_at"` // 修改时间
	DeletedAt     time.Time  `gorm:"column:deleted_at" json:"deleted_at"`                           // 删除时间
}

func (m *User) TableName() string {
	return "user"
}
