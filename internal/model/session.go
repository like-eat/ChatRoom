package model

import "time"

type Session struct {
	Id            int64     `gorm:"column:id;primaryKey;comment:自增id"`
	Uuid          string    `gorm:"column:uuid;uniqueIndex;type:char(20);comment:会话uuid"`
	LastMessage   string    `gorm:"column:last_message;type:varchar(2000);comment:最新的消息"`
	LastMessageAt time.Time `gorm:"column:last_message_at;type:datetime;comment:最近接收时间"`
}

func (Session) TableName() string {
	return "session"
}
