package gorm

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"kama_chat_server/internal/dao"
	"kama_chat_server/internal/dto/request"
	"kama_chat_server/internal/model"
	"kama_chat_server/pkg/util/random"
	"kama_chat_server/pkg/zlog"
)

type sessionService struct {
}

var SessionService = new(sessionService)

// CreateSession 创建会话
func (s *sessionService) CreateSession(req request.CreateSessionRequest) error {
	var user model.UserInfo
	if res := dao.GormDB.First(&user, req.SendId); res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			zlog.Error(" send user not found")
			return errors.New("user not found")
		} else {
			zlog.Error(res.Error.Error())
			return res.Error
		}
	}
	if user.DeletedAt.Valid {
		zlog.Info("user has been deleted")
		return nil
	}
	if req.ReceiveId[0] == 'U' {
		var receiveUser model.UserInfo
		if res := dao.GormDB.First(&receiveUser, req.ReceiveId); res.Error != nil {
			if errors.Is(res.Error, gorm.ErrRecordNotFound) {
				zlog.Error("receive user not found")
				return errors.New("receive user not found")
			} else {
				zlog.Error(res.Error.Error())
				return res.Error
			}
		}
		if receiveUser.DeletedAt.Valid {
			zlog.Info("receive user has been deleted")
			return nil
		}
	} else {
		var receiveGroup model.GroupInfo
		if res := dao.GormDB.First(&receiveGroup, req.ReceiveId); res.Error != nil {
			if errors.Is(res.Error, gorm.ErrRecordNotFound) {
				zlog.Error("receive group not found")
				return errors.New("receive group not found")
			} else {
				zlog.Error(res.Error.Error())
				return res.Error
			}
		}
		if receiveGroup.DeletedAt.Valid {
			zlog.Info("receive group has been deleted")
			return nil
		}
	}
	var session model.Session
	session.Uuid = fmt.Sprintf("S%d", random.GetNowAndLenRandomString(11))
	session.SendId = req.SendId
	session.ReceiveId = req.ReceiveId
	return dao.GormDB.Create(&session).Error
}
