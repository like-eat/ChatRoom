package dao

import (
	"apylee_chat_server/config"
	"apylee_chat_server/internal/model"
	"apylee_chat_server/pkg/log"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var gormClient *gorm.DB

func creatClient() (*gorm.DB, error) {
	if gormClient == nil {
		conf := config.GetConfig()
		user := conf.User
		password := conf.MysqlConfig.Password
		host := conf.MysqlConfig.Host
		port := conf.MysqlConfig.Port
		appName := conf.AppName
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, password, host, port, appName)
		fmt.Println(dsn)
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
