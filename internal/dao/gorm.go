package dao

import (
	"fmt"
	"kama_chat_server/internal/config"
	"kama_chat_server/internal/model"
	"kama_chat_server/pkg/zlog"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var GormDB *gorm.DB

// 在程序跑起来之前自动把数据库连接准备好
// 连接数据库并建立刘张表
func init() {
	// 获取MySQL，Redis，JWT，Kafka的配置
	conf := config.GetConfig()
	// 拼数据库连接串
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		conf.MysqlConfig.User,
		conf.MysqlConfig.Password,
		conf.MysqlConfig.Host,
		conf.MysqlConfig.Port,
		conf.MysqlConfig.DatabaseName,
	)

	var err error
	// 连接数据库
	GormDB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		zlog.Fatal(err.Error())
	}
	// 自动建表——刘张表
	err = GormDB.AutoMigrate(&model.UserInfo{}, &model.GroupInfo{}, &model.UserContact{}, &model.Session{}, &model.ContactApply{}, &model.Message{}) // 自动迁移，如果没有建表，会自动创建对应的表
	if err != nil {
		zlog.Fatal(err.Error())
	}
}
