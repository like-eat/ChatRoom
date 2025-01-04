package gorm

import (
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
	"kama_chat_server/internal/dao"
	"kama_chat_server/internal/dto/request"
	"kama_chat_server/internal/dto/respond"
	"kama_chat_server/internal/model"
	"kama_chat_server/pkg/constants"
	"kama_chat_server/pkg/enum/contact/contact_status_enum"
	"kama_chat_server/pkg/enum/contact/contact_type_enum"
	"kama_chat_server/pkg/enum/group_info/group_status_enum"
	"kama_chat_server/pkg/util/random"
	"kama_chat_server/pkg/zlog"
	"time"
)

type groupInfoService struct {
}

var GroupInfoService = new(groupInfoService)

// SaveGroup 保存群聊
//func (g *groupInfoService) SaveGroup(groupReq request.SaveGroupRequest) error {
//	var group model.GroupInfo
//	res := dao.GormDB.First(&group, "uuid = ?", groupReq.Uuid)
//	if res.Error != nil {
//		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
//			// 创建群聊
//			group.Uuid = groupReq.Uuid
//			group.Name = groupReq.Name
//			group.Notice = groupReq.Notice
//			group.AddMode = groupReq.AddMode
//			group.Avatar = groupReq.Avatar
//			group.MemberCnt = 1
//			group.Members = append(group.Members, groupReq.OwnerId)
//			group.OwnerId = groupReq.OwnerId
//			group.CreatedAt = time.Now()
//			group.UpdatedAt = time.Now()
//			if res := dao.GormDB.Create(&group); res.Error != nil {
//				zlog.Error(res.Error.Error())
//				return res.Error
//			}
//			return nil
//		} else {
//			zlog.Error(res.Error.Error())
//			return res.Error
//		}
//	}
//	// 更新群聊
//	group.Uuid = groupReq.Uuid
//	group.Name = groupReq.Name
//	group.Notice = groupReq.Notice
//	group.AddMode = groupReq.AddMode
//	group.Avatar = groupReq.Avatar
//	group.MemberCnt = 1
//	group.Members = append(group.Members, groupReq.OwnerId)
//	group.OwnerId = groupReq.OwnerId
//	group.CreatedAt = time.Now()
//	group.UpdatedAt = time.Now()
//	return nil
//}

// CreateGroup 创建群聊
func (g *groupInfoService) CreateGroup(groupReq request.CreateGroupRequest) (string, int) {
	group := model.GroupInfo{
		Uuid:      fmt.Sprintf("G%s", random.GetNowAndLenRandomString(11)),
		Name:      groupReq.Name,
		Notice:    groupReq.Notice,
		OwnerId:   groupReq.OwnerId,
		MemberCnt: 1,
		AddMode:   groupReq.AddMode,
		Avatar:    groupReq.Avatar,
		Status:    group_status_enum.NORMAL,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	var members []string
	members = append(members, groupReq.OwnerId)
	var err error
	group.Members, err = json.Marshal(members)
	if err != nil {
		zlog.Error(err.Error())
		return constants.SYSTEM_ERROR, -1
	}
	if res := dao.GormDB.Create(&group); res.Error != nil {
		zlog.Error(res.Error.Error())
		return constants.SYSTEM_ERROR, -1
	}
	return "创建成功", 0
}

// GetAllMembers 获取所有成员信息
//func (g *groupInfoService) GetAllMembers(groupId string) ([]string, int) {
//	var group model.GroupInfo
//	res := dao.GormDB.Preload("Members").First(&group, "uuid = ?", groupId)
//	if res.Error != nil {
//		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
//			zlog.Error("群组不存在")
//			return nil, -1
//		} else {
//			zlog.Error(res.Error.Error())
//			return nil, -1
//		}
//	} else {
//		var members []string
//		if err := json.Unmarshal(group.Members, members); err != nil {
//			zlog.Error(err.Error())
//			return nil, -1
//		}
//		return members, 0
//	}
//}

// LoadMyGroup 获取我创建的群聊
func (g *groupInfoService) LoadMyGroup(ownerId string) (string, []respond.LoadMyGroupRespond, int) {
	var groupList []model.GroupInfo
	if res := dao.GormDB.Order("created_at DESC").Where("owner_id = ?", ownerId).Find(&groupList); res.Error != nil {
		zlog.Error(res.Error.Error())
		return constants.SYSTEM_ERROR, nil, -1
	}
	var groupListRsp []respond.LoadMyGroupRespond
	for _, group := range groupList {
		groupListRsp = append(groupListRsp, respond.LoadMyGroupRespond{
			GroupId:   group.Uuid,
			GroupName: group.Name,
			Avatar:    group.Avatar,
		})
	}

	return "获取成功", groupListRsp, 0
}

// GetGroupInfo 获取聊天详情
func (g *groupInfoService) GetGroupInfo(userId string, groupId string) (string, *model.GroupInfo, int) {
	var user model.UserInfo
	if res := dao.GormDB.First(&user, "uuid = ?", userId); res.Error != nil {
		zlog.Error(res.Error.Error())
		return constants.SYSTEM_ERROR, nil, -1
	}
	var group model.GroupInfo
	if res := dao.GormDB.First(&group, "uuid = ?", groupId); res.Error != nil {
		zlog.Error(res.Error.Error())
		return constants.SYSTEM_ERROR, nil, -1
	}
	return "获取成功", &group, 0
}

//func (g *groupInfoService) checkUserAndGroupValid(userId string, groupId string)

// GetGroupInfo4Chat 获取聊天会话群聊详情
//func (g *groupInfoService) GetGroupInfo4Chat() error {
//
//}

// LeaveGroup 退群
func (g *groupInfoService) LeaveGroup(userId string, groupId string) (string, int) {
	// 从群聊中清除该用户
	var group model.GroupInfo
	if res := dao.GormDB.First(&group, "uuid = ?", groupId); res.Error != nil {
		zlog.Error(res.Error.Error())
		return constants.SYSTEM_ERROR, -1
	}
	var members []string
	if err := json.Unmarshal(group.Members, &members); err != nil {
		zlog.Error(err.Error())
		return constants.SYSTEM_ERROR, -1
	}
	for i, member := range members {
		if member == userId {
			members = append(members[:i], members[i+1:]...)
			break
		}
	}
	if data, err := json.Marshal(members); err != nil {
		zlog.Error(err.Error())
		return constants.SYSTEM_ERROR, -1
	} else {
		group.Members = data
	}
	group.MemberCnt -= 1
	if res := dao.GormDB.Save(&group); res.Error != nil {
		zlog.Error(res.Error.Error())
		return constants.SYSTEM_ERROR, -1
	}
	// 删除会话
	var deletedAt gorm.DeletedAt
	deletedAt.Time = time.Now()
	deletedAt.Valid = true
	if res := dao.GormDB.Model(&model.Session{}).Where("send_id = ? AND receive_id = ?", userId, groupId).Update("deleted_at", deletedAt); res.Error != nil {
		zlog.Error(res.Error.Error())
		return constants.SYSTEM_ERROR, -1
	}
	// 删除联系人
	if res := dao.GormDB.Model(&model.UserContact{}).Where("user_id = ? AND contact_id = ?", userId, groupId).Updates(map[string]interface{}{
		"deleted_at": deletedAt,
		"status":     contact_status_enum.QUIT_GROUP, // 退群
	}); res.Error != nil {
		zlog.Error(res.Error.Error())
		return constants.SYSTEM_ERROR, -1
	}
	// 删除申请记录，后面还可以加
	if res := dao.GormDB.Model(&model.ContactApply{}).Where("contact_id = ? AND user_id = ?", groupId, userId).Update("deleted_at", deletedAt); res.Error != nil {
		zlog.Error(res.Error.Error())
		return constants.SYSTEM_ERROR, -1
	}

	return "退群成功", 0
}

// DismissGroup 解散群聊
func (g *groupInfoService) DismissGroup(groupId string) (string, int) {
	var deletedAt gorm.DeletedAt
	deletedAt.Time = time.Now()
	deletedAt.Valid = true
	if res := dao.GormDB.Model(&model.GroupInfo{}).Where("uuid = ?", groupId).Updates(
		map[string]interface{}{
			"deleted_at": deletedAt,
			"updated_at": deletedAt.Time,
		}); res.Error != nil {
		zlog.Error(res.Error.Error())
		return constants.SYSTEM_ERROR, -1
	}

	var sessionList []model.Session
	if res := dao.GormDB.Model(&model.Session{}).Where("receive_id = ?", groupId).Find(&sessionList); res.Error != nil {
		zlog.Error(res.Error.Error())
		return constants.SYSTEM_ERROR, -1
	}
	for _, session := range sessionList {
		if res := dao.GormDB.Model(&session).Updates(
			map[string]interface{}{
				"deleted_at": deletedAt,
			}); res.Error != nil {
			zlog.Error(res.Error.Error())
			return constants.SYSTEM_ERROR, -1
		}
	}

	var userContactList []model.UserContact
	if res := dao.GormDB.Model(&model.UserContact{}).Where("contact_id = ?", groupId).Find(&userContactList); res.Error != nil {
		zlog.Error(res.Error.Error())
		return constants.SYSTEM_ERROR, -1
	}

	for _, userContact := range userContactList {
		if res := dao.GormDB.Model(&userContact).Update("deleted_at", deletedAt); res.Error != nil {
			zlog.Error(res.Error.Error())
			return constants.SYSTEM_ERROR, -1
		}
	}

	var contactApply model.ContactApply
	if res := dao.GormDB.Model(&contactApply).Where("contact_id = ?", groupId).Update("deleted_at", deletedAt); res.Error != nil {
		if res.Error != gorm.ErrRecordNotFound {
			zlog.Info(res.Error.Error())
			return "无响应的申请记录需要删除", 0
		}
		zlog.Error(res.Error.Error())
		return constants.SYSTEM_ERROR, -1
	}
	return "解散群聊成功", 0
}

// CheckGroupAddMode 检查群聊加群方式
func (g *groupInfoService) CheckGroupAddMode(groupId string) (string, int8, int) {
	var group model.GroupInfo
	if res := dao.GormDB.First(&group, "uuid = ?", groupId); res.Error != nil {
		zlog.Error(res.Error.Error())
		return constants.SYSTEM_ERROR, -1, -1
	}
	return "加群方式获取成功", group.AddMode, 0
}

// EnterGroupDirectly 直接进群
// ownerId 是群聊id
func (g *groupInfoService) EnterGroupDirectly(ownerId, contactId string) (string, int) {

	var group model.GroupInfo
	if res := dao.GormDB.First(&group, "uuid = ?", ownerId); res.Error != nil {
		zlog.Error(res.Error.Error())
		return constants.SYSTEM_ERROR, -1
	}
	var members []string
	if err := json.Unmarshal(group.Members, &members); err != nil {
		zlog.Error(err.Error())
		return constants.SYSTEM_ERROR, -1
	}
	members = append(members, contactId)
	if data, err := json.Marshal(members); err != nil {
		zlog.Error(err.Error())
		return constants.SYSTEM_ERROR, -1
	} else {
		group.Members = data
	}
	group.MemberCnt += 1
	if res := dao.GormDB.Save(&group); res.Error != nil {
		zlog.Error(res.Error.Error())
		return constants.SYSTEM_ERROR, -1
	}
	newContact := model.UserContact{
		UserId:      contactId,
		ContactId:   ownerId,
		ContactType: contact_type_enum.Group,    // 用户
		Status:      contact_status_enum.NORMAL, // 正常
		CreatedAt:   time.Now(),
		UpdateAt:    time.Now(),
	}
	if res := dao.GormDB.Create(&newContact); res.Error != nil {
		zlog.Error(res.Error.Error())
		return constants.SYSTEM_ERROR, -1
	}
	return "进群成功", 0
}
