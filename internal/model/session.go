package model

import "time"

type Session struct {
	Id            int64     `gorm:"column:id;primaryKey;comment:自增id"`
	Uuid          string    `gorm:"column:uuid;uniqueIndex;type:char(20);comment:会话uuid"`
	SendId        string    `gorm:"column:send_id;Index;type:char(20);not null;comment:创建会话人id"`
	ReceiveId     string    `gorm:"column:receive_id;type:char(20);not null;comment:接受会话人id"`
	LastMessage   string    `gorm:"column:last_message;type:varchar(2000);comment:最新的消息"`
	LastMessageAt time.Time `gorm:"column:last_message_at;type:datetime;comment:最近接收时间"`
}

func (Session) TableName() string {
	return "session"
}
