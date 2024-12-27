package gorm

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"kama_chat_server/internal/dao"
	"kama_chat_server/internal/dto/request"
	"kama_chat_server/internal/dto/respond"
	"kama_chat_server/internal/model"
	"kama_chat_server/pkg/enum/contact/contact_status_enum"
	"kama_chat_server/pkg/enum/error_info"
	"kama_chat_server/pkg/enum/group_info/group_status_enum"
	"kama_chat_server/pkg/enum/user_info/user_status_enum"
	"kama_chat_server/pkg/util/random"
	"kama_chat_server/pkg/zlog"
	"time"
)

type sessionService struct {
}

var SessionService = new(sessionService)

// CreateSession 创建会话
func (s *sessionService) CreateSession(req request.CreateSessionRequest) (string, string, int) {
	var user model.UserInfo
	if res := dao.GormDB.Where("uuid = ?", req.SendId).First(&user); res.Error != nil {
		zlog.Error(res.Error.Error())
		return error_info.SYSTEM_ERROR, "", -1
	}
	var session model.Session
	session.Uuid = fmt.Sprintf("S%s", random.GetNowAndLenRandomString(11))
	session.SendId = req.SendId
	session.ReceiveId = req.ReceiveId
	session.CreatedAt = time.Now()
	if req.ReceiveId[0] == 'U' {
		var receiveUser model.UserInfo
		if res := dao.GormDB.Where("uuid = ?", req.ReceiveId).First(&receiveUser); res.Error != nil {
			zlog.Error(res.Error.Error())
			return error_info.SYSTEM_ERROR, "", -1
		}
		if receiveUser.Status == user_status_enum.DISABLE {
			zlog.Error("该用户被禁用了")
			return "该用户被禁用了", "", -2
		} else {
			session.ReceiveName = receiveUser.Nickname
			session.Avatar = receiveUser.Avatar
		}
	} else {
		var receiveGroup model.GroupInfo
		if res := dao.GormDB.Where("uuid = ?", req.ReceiveId).First(&receiveGroup); res.Error != nil {
			zlog.Error(res.Error.Error())
			return error_info.SYSTEM_ERROR, "", -1
		}
		if receiveGroup.Status == group_status_enum.DISABLE {
			zlog.Error("该群聊被禁用了")
			return "该群聊被禁用了", "", -2
		} else {
			session.ReceiveName = receiveGroup.Name
			session.Avatar = receiveGroup.Avatar
		}
	}

	if res := dao.GormDB.Create(&session); res.Error != nil {
		zlog.Error(res.Error.Error())
		return error_info.SYSTEM_ERROR, "", -1
	}
	return "会话创建成功", session.Uuid, 0
}

// CheckOpenSessionAllowed 检查是否允许发起会话
func (s *sessionService) CheckOpenSessionAllowed(sendId, receiveId string) (string, bool, int) {
	var contact model.UserContact
	if res := dao.GormDB.Where("user_id = ? and contact_id = ?", sendId, receiveId).First(&contact); res.Error != nil {
		zlog.Error(res.Error.Error())
		return error_info.SYSTEM_ERROR, false, -1
	}
	if contact.Status == contact_status_enum.BE_BLACK {
		return "已被对方拉黑，无法发起会话", false, 0
	} else if contact.Status == contact_status_enum.BLACK {
		return "已拉黑对方，先解除拉黑状态才能发起会话", false, 0
	}
	return "可以发起会话", true, 0
}

// DeleteSession 删除会话

// OpenSession 打开会话
func (s *sessionService) OpenSession(req request.OpenSessionRequest) (string, string, int) {
	var session model.Session
	if res := dao.GormDB.Where("send_id = ? and receive_id = ?", req.SendId, req.ReceiveId).First(&session); res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			zlog.Info("会话没有找到，将新建会话")
			createReq := request.CreateSessionRequest{
				SendId:    req.SendId,
				ReceiveId: req.ReceiveId,
			}
			return s.CreateSession(createReq)
		}
	}
	return "会话创建成功", session.Uuid, 0
}

// GetUserSessionList 获取用户会话列表
func (s *sessionService) GetUserSessionList(ownerId string) (string, []respond.UserSessionListRespond, int) {
	var sessionList []model.Session
	if res := dao.GormDB.Order("created_at DESC").Where("send_id = ?", ownerId).Find(&sessionList); res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			zlog.Info("未创建用户会话")
			return "未创建用户会话", nil, 0
		} else {
			zlog.Error(res.Error.Error())
			return error_info.SYSTEM_ERROR, nil, -1
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
	return "获取成功", sessionListRsp, 0
}

// GetGroupSessionList 获取群聊会话列表
func (s *sessionService) GetGroupSessionList(ownerId string) (string, []respond.GroupSessionListRespond, int) {
	var sessionList []model.Session
	if res := dao.GormDB.Order("created_at DESC").Where("send_id = ?", ownerId).Find(&sessionList); res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			zlog.Info("未创建群聊会话")
			return "未创建群聊会话", nil, 0
		} else {
			zlog.Error(res.Error.Error())
			return error_info.SYSTEM_ERROR, nil, -1
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
	return "获取成功", sessionListRsp, 0
}

// DeleteSession 删除会话
func (s *sessionService) DeleteSession(sessionId string) (string, int) {
	var session model.Session
	if res := dao.GormDB.Where("uuid = ?", sessionId).Find(&session); res.Error != nil {
		zlog.Error(res.Error.Error())
		return error_info.SYSTEM_ERROR, -1
	}
	session.DeletedAt.Valid = true
	session.DeletedAt.Time = time.Now()
	if res := dao.GormDB.Save(&session); res.Error != nil {
		zlog.Error(res.Error.Error())
		return error_info.SYSTEM_ERROR, -1
	}
	return "删除成功", 0
}
