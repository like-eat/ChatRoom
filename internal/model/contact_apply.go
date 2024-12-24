package model

import "time"

type ContactApply struct {
	Id          int64     `gorm:"column:id;primaryKey;comment:自增id"`
	Uuid        string    `gorm:"column:uuid;uniqueIndex;type:char(20);comment:申请id"`
	UserId      string    `gorm:"column:user_id;index;type:char(20);not null;comment:申请人id"`
	ContactId   string    `gorm:"column:contact_id;index;type:char(20);not null;comment:被申请id"`
	ContactType bool      `gorm:"column:contact_type;type:tinyint(1);not null;comment:被申请类型，0.用户，1.群聊"`
	Status      bool      `gorm:"column:status;not null;type:tinyint(1);comment:申请状态，0.申请中，1.申请通过"`
	Message     string    `gorm:"column:message;type:varchar(100);comment:申请信息"`
	LastApplyAt time.Time `gorm:"column:last_apply_at;type:datetime;not null;comment:最后申请时间"`
}

func (ContactApply) TableName() string {
	return "user_contact_apply"
}
