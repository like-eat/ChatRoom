package model

import (
	"gorm.io/gorm"
	"time"
)

type UserInfo struct {
	Id            int64          `gorm:"column:id;primaryKey;comment:自增id"`
	Uuid          string         `gorm:"column:uuid;uniqueIndex;type:varchar(12);comment:用户唯一id"`
	NickName      string         `gorm:"column:nickname;type:varchar(20);not null;comment:昵称"`
	TelePhone     string         `gorm:"column:telephone;uniqueIndex;not null;type:varchar(11);comment:电话"`
	Email         string         `gorm:"column:email;type:varchar(30);comment:邮箱"`
	Avatar        string         `gorm:"column:avatar;type:varchar(255);default:default_avatar.png;not null;comment:头像"`
	Gender        bool           `gorm:"column:gender;type:tinyint(1);comment:性别"`
	Signature     string         `gorm:"column:signature;type:varchar(50);comment:个性签名"`
	Password      string         `gorm:"column:password;type:char(18);not null;comment:密码"`
	CreatedAt     time.Time      `gorm:"column:created_at;type:datetime;not null;comment:创建时间"`
	DeletedAt     gorm.DeletedAt `gorm:"column:deleted_at;type:datetime;comment:删除时间"`
	LastOnlineAt  time.Time      `gorm:"column:last_online_at;type:datetime;comment:上次登录时间"`
	LastOfflineAt time.Time      `gorm:"column:last_offline_at;type:datetime;comment:最近离线时间"`
	IsAdmin       bool           `gorm:"column:is_admin;type:tinyint(1);not null;comment:是否是管理员，0.不是，1.是"`
}

func (UserInfo) TableName() string {
	return "user_info"
}
