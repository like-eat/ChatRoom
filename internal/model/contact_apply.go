package model

import (
	"time"

	"gorm.io/gorm"
)

// 好友申请表 记录加好友加群的申请记录
type ContactApply struct {
	Id          int64          `gorm:"column:id;primaryKey;comment:自增id"`
	Uuid        string         `gorm:"column:uuid;uniqueIndex;type:char(20);comment:申请id"`
	UserId      string         `gorm:"column:user_id;index;type:char(20);not null;comment:申请人id"`
	ContactId   string         `gorm:"column:contact_id;index;type:char(20);not null;comment:被申请id"`
	ContactType int8           `gorm:"column:contact_type;not null;comment:被申请类型，0.用户，1.群聊"`
	Status      int8           `gorm:"column:status;not null;comment:申请状态，0.申请中，1.通过，2.拒绝"`
	Message     string         `gorm:"column:message;type:varchar(100);comment:申请信息"`
	LastApplyAt time.Time      `gorm:"column:last_apply_at;type:datetime;not null;comment:最后申请时间"`
	DeletedAt   gorm.DeletedAt `gorm:"column:deleted_at;index;type:datetime;comment:删除时间"`
}

func (ContactApply) TableName() string {
	return "contact_apply"
}
