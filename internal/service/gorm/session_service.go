package gorm

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"kama_chat_server/internal/dao"
	"kama_chat_server/internal/dto/request"
	"kama_chat_server/internal/dto/respond"
	"kama_chat_server/internal/model"
	"kama_chat_server/pkg/util/random"
	"kama_chat_server/pkg/zlog"
	"time"
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
	session.CreatedAt = time.Now()
	if req.ReceiveId[0] == 'U' {
		var user model.UserInfo
		if res := dao.GormDB.First(&user, req.ReceiveId); res.Error != nil {
			if errors.Is(res.Error, gorm.ErrRecordNotFound) {
				zlog.Error("receive user not found")
				return errors.New("receive user not found")
			} else {
				zlog.Error(res.Error.Error())
				return res.Error
			}
		}
		if !user.DeletedAt.Valid {
			session.ReceiveName = user.Nickname
			session.Avatar = user.Avatar
		}
	} else {
		var group model.GroupInfo
		if res := dao.GormDB.First(&group, req.ReceiveId); res.Error != nil {
			if errors.Is(res.Error, gorm.ErrRecordNotFound) {
				zlog.Error("receive group not found")
				return errors.New("receive group not found")
			} else {
				zlog.Error(res.Error.Error())
				return res.Error
			}
		}
		if !group.DeletedAt.Valid {
			session.ReceiveName = group.Name
			session.Avatar = group.Avatar
		}
	}
	return dao.GormDB.Create(&session).Error
}

// OpenSession 打开会话
func (s *sessionService) OpenSession(req request.OpenSessionRequest) error {
	var session model.Session
	if res := dao.GormDB.Where("send_id = ? and receive_id = ?", req.SendId, req.ReceiveId).First(&session); res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			zlog.Info("session not found")
			createReq := request.CreateSessionRequest{
				SendId:    req.SendId,
				ReceiveId: req.ReceiveId,
			}
			if err := s.CreateSession(createReq); err != nil {
				return err
			}
		}
	}
	return nil
}

// GetUserSessionList 获取用户会话列表
func (s *sessionService) GetUserSessionList(ownerId string) ([]respond.UserSessionListRespond, error) {
	var sessionList []model.Session
	if res := dao.GormDB.Order("created_at DESC").Where("send_id = ?", ownerId).Find(&sessionList); res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			zlog.Info("session not found")
			return nil, nil
		} else {
			zlog.Error(res.Error.Error())
			return nil, res.Error
		}
	}
	var sessionListRsp []respond.UserSessionListRespond
	for i := 0; i < len(sessionList); i++ {
		if sessionList[i].ReceiveId[0] == 'U' {
			sessionListRsp = append(sessionListRsp, respond.UserSessionListRespond{
				SessionId: sessionList[i].Uuid,
				Avatar:    sessionList[i].Avatar,
				UserId:    sessionList[i].ReceiveId,
				Username:  sessionList[i].ReceiveName,
			})
		}
	}
	return sessionListRsp, nil
}

// GetGroupSessionList 获取群聊会话列表
func (s *sessionService) GetGroupSessionList(ownerId string) ([]respond.GroupSessionListRespond, error) {
	var sessionList []model.Session
	if res := dao.GormDB.Order("created_at DESC").Where("send_id = ?", ownerId).Find(&sessionList); res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			zlog.Info("session not found")
			return nil, nil
		} else {
			zlog.Error(res.Error.Error())
			return nil, res.Error
		}
	}
	var sessionListRsp []respond.GroupSessionListRespond
	for i := 0; i < len(sessionList); i++ {
		if sessionList[i].ReceiveId[0] == 'G' {
			sessionListRsp = append(sessionListRsp, respond.GroupSessionListRespond{
				SessionId: sessionList[i].Uuid,
				Avatar:    sessionList[i].Avatar,
				GroupId:   sessionList[i].ReceiveId,
				GroupName: sessionList[i].ReceiveName,
			})
		}
	}
	return sessionListRsp, nil
}
