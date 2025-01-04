package gorm

import (
	"fmt"
	"kama_chat_server/internal/dao"
	"kama_chat_server/internal/dto/respond"
	"kama_chat_server/internal/model"
	"kama_chat_server/pkg/constants"
	"kama_chat_server/pkg/zlog"
)

type messageService struct {
}

var MessageService = new(messageService)

// GetMessageList 获取聊天记录
func (m *messageService) GetMessageList(userOneId, userTwoId string) (string, []respond.GetMessageListRespond, int) {
	zlog.Info(fmt.Sprintf("%s %s", userTwoId, userTwoId))
	var messageList []model.Message
	if res := dao.GormDB.Where("(send_id = ? AND receive_id = ?) OR (send_id = ? AND receive_id = ?)", userOneId, userTwoId, userTwoId, userOneId).Order("created_at ASC").Find(&messageList); res.Error != nil {
		zlog.Error(res.Error.Error())
		return constants.SYSTEM_ERROR, nil, -1
	}
	var rspList []respond.GetMessageListRespond
	for _, message := range messageList {
		rspList = append(rspList, respond.GetMessageListRespond{
			SendId:    message.SendId,
			SendName:  message.SendName,
			ReceiveId: message.ReceiveId,
			Content:   message.Content,
			Url:       message.Url,
			Type:      message.Type,
			FileType:  message.FileType,
			FileName:  message.FileName,
			FileSize:  message.FileSize,
			CreatedAt: message.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}
	return "获取聊天记录成功", rspList, 0
}

// GetGroupMessageList 获取群聊消息记录
func (m *messageService) GetGroupMessageList(groupId string) (string, []respond.GetMessageListRespond, int) {
	var messageList []model.Message
	if res := dao.GormDB.Where("receive_id = ?", groupId).Order("created_at ASC").Find(&messageList); res.Error != nil {
		zlog.Error(res.Error.Error())
		return constants.SYSTEM_ERROR, nil, -1
	}
	var rspList []respond.GetMessageListRespond
	for _, message := range messageList {
		rsp := respond.GetMessageListRespond{
			SendId:    message.SendId,
			SendName:  message.SendName,
			ReceiveId: message.ReceiveId,
			Content:   message.Content,
			Url:       message.Url,
			Type:      message.Type,
			FileType:  message.FileType,
			FileName:  message.FileName,
			FileSize:  message.FileSize,
			CreatedAt: message.CreatedAt.Format("2006-01-02 15:04:05"),
		}
		rspList = append(rspList, rsp)
	}
	return "获取聊天记录成功", rspList, 0
}
