package model

import "time"

type Message struct {
	Id        int64     `gorm:"column:id;primaryKey;comment:自增id"`
	Uuid      string    `gorm:"column:uuid;uniqueIndex;type:char(20);not null;comment:消息uuid"`
	SessionId string    `gorm:"column:session_id;uniqueIndex;type:char(20);not null;comment:会话uuid"`
	Type      int8      `gorm:"column:type;not null;comment:消息类型，0.文本，1.语音，2.文件，3.通话"`
	Content   string    `gorm:"column:content;type:char(10000);comment:消息内容"`
	Url       string    `gorm:"column:url;type:char(100);comment:消息url"`
	SendId    string    `gorm:"column:send_id;type:char(20);not null;comment:发送者uuid"`
	ReceiveId string    `gorm:"column:receive_id;type:char(20);not null;comment:接受者uuid"`
	FileType  string    `gorm:"column:file_type;type:char(10);not null;comment:文件类型"`
	FileName  string    `gorm:"column:file_name;type:varchar(50);not null;comment:文件名"`
	FileSize  int       `gorm:"column:file_size;not null;comment:文件大小"`
	Status    int8      `gorm:"column:status;not null;comment:状态，0.未发送，1.已发送"`
	CreatedAt time.Time `gorm:"column:created_at;not null;comment:创建时间"`
}

func (Message) TableName() string {
	return "message"
}
