package dao

import (
	"kama_chat_server/internal/dao"
	"kama_chat_server/internal/model"
	"kama_chat_server/pkg/util/random"
	"strconv"
	"testing"
	"time"
)

func TestCreate(t *testing.T) {
	userInfo := &model.UserInfo{
		Uuid:      "U" + strconv.Itoa(random.GetRandomInt(11)),
		Nickname:  "apylee",
		Telephone: "1" + strconv.Itoa(random.GetRandomInt(10)),
		Email:     "1212312312@qq.com",
		Password:  "123456",
		CreatedAt: time.Now(),
		IsAdmin:   1,
	}
	if err := dao.GormDB.Create(userInfo).Error; err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := dao.GormDB.Unscoped().Where("uuid = ?", userInfo.Uuid).Delete(&model.UserInfo{}).Error; err != nil {
			t.Error(err)
		}
	})
}
