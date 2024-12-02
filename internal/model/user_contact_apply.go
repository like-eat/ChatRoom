package model

import "time"

type UserContactApply struct {
	Id          int64     `gorm:"column:id;primaryKey;comment:申请id"`
	UserId      string    `gorm:"column:user_id;index;type:char(20);not null;comment:申请人id"`
	ContactId   string    `gorm:"column:apply_contact_id;index;type:char(20);not null;comment:被申请人id"`
	ContactType bool      `gorm:"column:apply_contact_type;type:tinyint(1);not null;comment:被申请人类型"`
	Status      bool      `gorm:"column:status;not null;type:tinyint(1);comment:申请状态，0.申请中，1.申请通过"`
	Message     string    `gorm:"column:message;type:varchar(100);comment:申请信息"`
	LastApplyAt time.Time `gorm:"column:last_apply_time;type:datetime;not null;comment:最后申请时间"`
}

func (UserContactApply) TableName() string {
	return "user_contact_apply"
}
