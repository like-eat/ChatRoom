package gorm

import (
	"errors"
	"gorm.io/gorm"
	"kama_chat_server/internal/dao"
	"kama_chat_server/internal/dto/respond"
	"kama_chat_server/internal/model"
	"kama_chat_server/pkg/zlog"
	"log"
)

type userContactService struct {
}

var UserContactService = new(userContactService)

func (u *userContactService) GetUserList(ownerId string) (string, []respond.MyUserListRespond, error) {
	var userList []model.UserContact
	if res := dao.GormDB.Order("created_at DESC").Where("user_id = ?", ownerId).Find(&userList); res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			message := "contact not found"
			zlog.Info(message)
			return message, nil, nil
		} else {
			zlog.Error(res.Error.Error())
			return "", nil, res.Error
		}
	}
	log.Println(userList)
	var userListRsp []respond.MyUserListRespond
	for _, user := range userList {
		contactId := user.ContactId
		var user model.UserInfo
		if res := dao.GormDB.First(&user, "uuid = ?", contactId); res.Error != nil {
			if errors.Is(res.Error, gorm.ErrRecordNotFound) {
				zlog.Error("contact not found in UserInfo table")
				return "", nil, res.Error
			} else {
				zlog.Error(res.Error.Error())
				return "", nil, res.Error
			}
		}
		userListRsp = append(userListRsp, respond.MyUserListRespond{
			UserId:   user.Uuid,
			UserName: user.Nickname,
			Avatar:   user.Avatar,
		})
	}
	//userListStr, err := json.Marshal(userListRsp)
	//if err != nil {
	//	return "", userListRsp, err
	//}
	return "", userListRsp, nil
}
