package model

import (
	"gorm.io/gorm"
	"time"
)

type GroupInfo struct {
	Id        string `gorm:"primaryKey"`
	Name      string `gorm:"size:20"`
	Notice    string `gorm:"size:500"`
	MemberCnt int
	OwnerId   string
	AddMode   int // 加群方式
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
