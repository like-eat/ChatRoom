package dao

import (
	"apylee_chat_server/internal/dao"
	"apylee_chat_server/internal/model"
	"apylee_chat_server/pkg/util/random"
	"strconv"
	"testing"
	"time"
)

func TestCreate(t *testing.T) {
	userInfo := &model.UserInfo{
		Uuid:      "U" + strconv.Itoa(random.GetRandomInt(11)),
		NickName:  "apylee",
		TelePhone: "180323532112",
		Email:     "1212312312@qq.com",
		Password:  "123456",
		CreatedAt: time.Now(),
		IsAdmin:   true,
	}
	_ = dao.GormDB.Create(userInfo)
}
