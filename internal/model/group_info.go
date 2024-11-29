package model

import (
	"gorm.io/gorm"
	"time"
)

type GroupInfo struct {
	Id        string `dao:"primaryKey"`
	Name      string `dao:"size:20"`
	Notice    string `dao:"size:500"`
	MemberCnt int
	OwnerId   string
	AddMode   int // 加群方式
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `dao:"index"`
}
