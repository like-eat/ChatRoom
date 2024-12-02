package gorm

import (
	"apylee_chat_server/internal/dao"
	"apylee_chat_server/internal/model"
	"apylee_chat_server/pkg/zlog"
	"errors"
	"gorm.io/gorm"
)

type groupInfoService struct {
}

var GroupInfoService = new(groupInfoService)

func (g *groupInfoService) SaveGroup() {

}

// GetAllMembers 获取所有成员信息
func (g *groupInfoService) GetAllMembers(groupId string) ([]model.UserInfo, error) {
	var group model.GroupInfo
	res := dao.GormDB.Preload("Members").First(&group, "uuid = ?", groupId)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			zlog.Info("group not existed")
			return nil, nil
		} else {
			zlog.Error(res.Error.Error())
			return nil, res.Error
		}
	} else {
		return group.Members, nil
	}
}
