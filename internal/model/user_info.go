package model

import (
	"gorm.io/gorm"
	"time"
)

type UserInfo struct {
	Id        string `dao:"primaryKey"`
	NickName  string `dao:"size:20"`
	TelePhone string `dao:"uniqueIndex;not null"`
	Email     string
	Avatar    string
	Age       int
	Gender    bool
	Signature string `dao:"size:50"` // 个性签名
	Password  string `dao:"size:20"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `dao:"index"`
}

var UserTable []UserInfo
