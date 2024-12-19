package gorm

import (
	"errors"
	"gorm.io/gorm"
	"kama_chat_server/internal/dao"
	"kama_chat_server/internal/dto/respond"
	"kama_chat_server/internal/model"
	"kama_chat_server/pkg/zlog"
)

type userContactService struct {
}

var UserContactService = new(userContactService)

// GetUserList 获取用户列表
func (u *userContactService) GetUserList(ownerId string) (string, []respond.MyUserListRespond, error) {
	var contactList []model.UserContact
	if res := dao.GormDB.Order("created_at DESC").Where("user_id = ?", ownerId).Find(&contactList); res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			message := "contact not found"
			zlog.Info(message)
			return message, nil, nil
		} else {
			zlog.Error(res.Error.Error())
			return "", nil, res.Error
		}
	}
	var userListRsp []respond.MyUserListRespond
	for _, contact := range contactList {
		if contact.ContactType == false {
			var user model.UserInfo
			if res := dao.GormDB.First(&user, "uuid = ?", contact.ContactId); res.Error != nil {
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
	}
	return "", userListRsp, nil
}

// LoadMyJoinedGroup 获取我加入的群聊
func (u *userContactService) LoadMyJoinedGroup(ownerId string) ([]respond.LoadMyJoinedGroupRespond, error) {
	var contactList []model.UserContact
	if res := dao.GormDB.Order("created_at DESC").Where("user_id = ?", ownerId).Find(&contactList); res.Error != nil {
		zlog.Error(res.Error.Error())
		return nil, res.Error
	}
	var groupList []model.GroupInfo
	for _, contact := range contactList {
		if contact.ContactId[0] == 'G' {
			var group model.GroupInfo
			if res := dao.GormDB.First(&group, "uuid = ?", contact.ContactId); res.Error != nil {
				if errors.Is(res.Error, gorm.ErrRecordNotFound) {
					zlog.Error("contact not found in GroupInfo")
					return nil, res.Error
				} else {
					zlog.Error(res.Error.Error())
					return nil, res.Error
				}
			}
			// 群没被删除，同时群主不是自己
			if !group.DeletedAt.Valid && group.OwnerId != ownerId {
				groupList = append(groupList, group)
			}
		}
	}
	var groupListRsp []respond.LoadMyJoinedGroupRespond
	for _, group := range groupList {
		groupListRsp = append(groupListRsp, respond.LoadMyJoinedGroupRespond{
			GroupId:   group.Uuid,
			GroupName: group.Name,
			Avatar:    group.Avatar,
		})
	}

	return groupListRsp, nil

}

// GetContactInfo 获取联系人信息
func (u *userContactService) GetContactInfo(contactId string) (string, *respond.GetContactInfoRespond, error) {
	if contactId[0] == 'G' {
		var group model.GroupInfo
		if res := dao.GormDB.First(&group, "uuid = ?", contactId); res.Error != nil {
			if errors.Is(res.Error, gorm.ErrRecordNotFound) {
				zlog.Error("contact not found in GroupInfo")
				return "", nil, res.Error
			} else {
				zlog.Error(res.Error.Error())
				return "", nil, res.Error
			}
		}
		if !group.DeletedAt.Valid {
			return "", &respond.GetContactInfoRespond{
				ContactId:        group.Uuid,
				ContactName:      group.Name,
				ContactAvatar:    group.Avatar,
				ContactNotice:    group.Notice,
				ContactAddMode:   group.AddMode,
				ContactMembers:   group.Members,
				ContactMemberCnt: group.MemberCnt,
				ContactOwnerId:   group.OwnerId,
			}, nil
		} else {
			zlog.Info("contact has been deleted")
			return "contact has been deleted", nil, nil
		}
	} else {
		var user model.UserInfo
		if res := dao.GormDB.First(&user, "uuid = ?", contactId); res.Error != nil {
			if errors.Is(res.Error, gorm.ErrRecordNotFound) {
				zlog.Error("contact not found in UserInfo")
				return "", nil, res.Error
			} else {
				zlog.Error(res.Error.Error())
				return "", nil, res.Error
			}
		}
		if !user.DeletedAt.Valid {
			return "", &respond.GetContactInfoRespond{
				ContactId:        user.Uuid,
				ContactName:      user.Nickname,
				ContactAvatar:    user.Avatar,
				ContactBirthday:  user.Birthday,
				ContactEmail:     user.Email,
				ContactPhone:     user.Telephone,
				ContactGender:    user.Gender,
				ContactSignature: user.Signature,
			}, nil
		} else {
			zlog.Info("contact has been deleted")
			return "contact has been deleted", nil, nil
		}
	}
}

// DeleteContact 删除联系人
func (u *userContactService) DeleteContact(ownerId, contactId string) error {
	if res := dao.GormDB.Where("user_id = ? AND contact_id = ?", ownerId, contactId).Delete(&model.UserContact{}); res.Error != nil {
		zlog.Error(res.Error.Error())
		return res.Error
	}
	if res := dao.GormDB.Where("send_id = ? AND receive_id = ?", ownerId, contactId).Delete(&model.Session{}); res.Error != nil {
		zlog.Error(res.Error.Error())
		return res.Error
	}
	return nil
}
