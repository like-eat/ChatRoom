package dao

import (
	"apylee_chat_server/internal/model"
	"apylee_chat_server/pkg/log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var gormClient *gorm.DB

func creatClient() (*gorm.DB, error) {
	if gormClient == nil {
		dsn := "root:dancernumber1@tcp(127.0.0.1:3306)/apylee-chat-server?charset=utf8mb4&parseTime=True&loc=Local"
		var err error
		gormClient, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err != nil {
			log.GetZapLogger().Fatal("failed to connect database")
			return nil, err
		}
		err = gormClient.AutoMigrate(&model.UserInfo{})
		if err != nil {
			log.GetZapLogger().Fatal("failed to migrate")
			return nil, err
		}
	}
	return gormClient, nil
}
