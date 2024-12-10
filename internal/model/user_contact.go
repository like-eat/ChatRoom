package model

import "time"

type UserContact struct {
	Id          int64     `gorm:"column:id;primaryKey;comment:自增id"`
	UserId      string    `gorm:"column:user_id;index;type:char(20);not null;comment:用户唯一id"`
	ContactId   string    `gorm:"column:contact_id;type:char(20);not null;comment:对应联系id"`
	ContactType bool      `gorm:"column:contact_type;type:tinyint(1);not null;comment:联系类型，0.用户，1.群聊"`
	Status      int       `gorm:"column:status;not null;comment:联系状态，0.正常，1.拉黑，2.被拉黑，3.删除好友，4.被删除好友"`
	CreatedAt   time.Time `gorm:"column:created_at;type:datetime;not null;comment:创建时间"`
	UpdateAt    time.Time `gorm:"column:update_at;type:datetime;not null;comment:更新时间"`
}

func (UserContact) TableName() string {
	return "user_contact"
}
