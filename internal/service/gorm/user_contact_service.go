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
	"kama_chat_server/pkg/enum/contact/contact_type_enum"
	"kama_chat_server/pkg/enum/contact_apply/contact_apply_status_enum"
	"kama_chat_server/pkg/enum/error_info"
	"kama_chat_server/pkg/enum/group_info/group_status_enum"
	"kama_chat_server/pkg/enum/user_info/user_status_enum"
	"kama_chat_server/pkg/util/random"
	"kama_chat_server/pkg/zlog"
	"time"
)

type userContactService struct {
}

var UserContactService = new(userContactService)

// GetUserList 获取用户列表
// 关于用户被禁用的问题，这里查到的是所有联系人，如果被禁用或被拉黑会以弹窗的形式提醒，无法打开会话框；如果被删除，是搜索不到该联系人的。
func (u *userContactService) GetUserList(ownerId string) (string, []respond.MyUserListRespond) {
	// dao
	var contactList []model.UserContact
	// 没有被删除
	if res := dao.GormDB.Order("created_at DESC").Where("user_id = ? AND status != 4", ownerId).Find(&contactList); res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			message := "目前不存在联系人"
			zlog.Error(message)
			return message, nil
		} else {
			zlog.Error(res.Error.Error())
			return error_info.SYSTEM_ERROR, nil
		}
	}
	// dto
	var userListRsp []respond.MyUserListRespond
	for _, contact := range contactList {
		// 联系人中是用户的
		if contact.ContactType == contact_type_enum.User {
			// 获取用户信息
			var user model.UserInfo
			if res := dao.GormDB.First(&user, "uuid = ?", contact.ContactId); res.Error != nil {
				zlog.Error(res.Error.Error())
				return error_info.SYSTEM_ERROR, nil
			}
			userListRsp = append(userListRsp, respond.MyUserListRespond{
				UserId:   user.Uuid,
				UserName: user.Nickname,
				Avatar:   user.Avatar,
			})
		}
	}
	return "获取用户列表成功", userListRsp
}

// LoadMyJoinedGroup 获取我加入的群聊
func (u *userContactService) LoadMyJoinedGroup(ownerId string) (string, []respond.LoadMyJoinedGroupRespond) {
	var contactList []model.UserContact
	// 没有退群，也没有被踢出群聊
	if res := dao.GormDB.Order("created_at DESC").Where("user_id = ? AND status != 6 AND status != 7", ownerId).Find(&contactList); res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			message := "目前不存在加入的群聊"
			zlog.Error(message)
			return message, nil
		} else {
			zlog.Error(res.Error.Error())
			return error_info.SYSTEM_ERROR, nil
		}
	}
	var groupList []model.GroupInfo
	for _, contact := range contactList {
		if contact.ContactId[0] == 'G' {
			var group model.GroupInfo
			if res := dao.GormDB.First(&group, "uuid = ?", contact.ContactId); res.Error != nil {
				zlog.Error(res.Error.Error())
				return error_info.SYSTEM_ERROR, nil
			}
			// 群没被删除，同时群主不是自己
			// 群主删除或admin删除群聊，status为7，即被踢出群聊，所以不用判断群是否被删除，删除了到不了这步
			if group.OwnerId != ownerId {
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

	return "获取加入群成功", groupListRsp
}

// GetContactInfo 获取联系人信息
// 调用这个接口的前提是该联系人没有处在删除或被删除，或者该用户还在群聊中
func (u *userContactService) GetContactInfo(contactId string) (string, *respond.GetContactInfoRespond) {
	if contactId[0] == 'G' {
		var group model.GroupInfo
		if res := dao.GormDB.First(&group, "uuid = ?", contactId); res.Error != nil {
			zlog.Error(res.Error.Error())
			return error_info.SYSTEM_ERROR, nil
		}
		// 没被禁用
		if group.Status != group_status_enum.DISABLE {
			return "", &respond.GetContactInfoRespond{
				ContactId:        group.Uuid,
				ContactName:      group.Name,
				ContactAvatar:    group.Avatar,
				ContactNotice:    group.Notice,
				ContactAddMode:   group.AddMode,
				ContactMembers:   group.Members,
				ContactMemberCnt: group.MemberCnt,
				ContactOwnerId:   group.OwnerId,
			}
		} else {
			zlog.Error("该群聊处于禁用状态")
			return "该群聊处于禁用状态", nil
		}
	} else {
		var user model.UserInfo
		if res := dao.GormDB.First(&user, "uuid = ?", contactId); res.Error != nil {
			zlog.Error(res.Error.Error())
			return error_info.SYSTEM_ERROR, nil
		}
		if user.Status != user_status_enum.DISABLE {
			return "", &respond.GetContactInfoRespond{
				ContactId:        user.Uuid,
				ContactName:      user.Nickname,
				ContactAvatar:    user.Avatar,
				ContactBirthday:  user.Birthday,
				ContactEmail:     user.Email,
				ContactPhone:     user.Telephone,
				ContactGender:    user.Gender,
				ContactSignature: user.Signature,
			}
		} else {
			zlog.Info("该用户处于禁用状态")
			return "该用户处于禁用状态", nil
		}
	}
}

// DeleteContact 删除联系人
func (u *userContactService) DeleteContact(ownerId, contactId string) {
	if res := dao.GormDB.Model(&model.UserContact{}).Where("user_id = ? AND contact_id = ?", ownerId, contactId).Update("deleted_at", time.Now()); res.Error != nil {
		zlog.Error(res.Error.Error())
		return
	}

	var userContact model.UserContact
	if res := dao.GormDB.Where("user_id = ? AND contact_id = ?", ownerId, contactId).Find(&userContact); res.Error != nil {
		zlog.Error(res.Error.Error())
		return
	}
	userContact.Status = contact_status_enum.DELETE
	userContact.DeletedAt = time.Now()

	if res := dao.GormDB.Model(&model.Session{}).Where("send_id = ? AND receive_id = ?", ownerId, contactId).Update("deleted_at", time.Now()); res.Error != nil {
		zlog.Error(res.Error.Error())
		return
	}

}

// ApplyContact 申请添加联系人
func (u *userContactService) ApplyContact(req request.ApplyContactRequest) (string, error) {
	if req.ContactId[0] == 'U' {
		var user model.UserInfo
		if res := dao.GormDB.First(&user, "uuid = ?", req.ContactId); res.Error != nil {
			if errors.Is(res.Error, gorm.ErrRecordNotFound) {
				zlog.Info("contact not found in UserInfo")
				return "用户不存在/已被禁用", nil
			} else {
				zlog.Error(res.Error.Error())
				return "", res.Error
			}
		}
		var contactApply model.ContactApply
		if res := dao.GormDB.Where("user_id = ? AND contact_id = ?", req.OwnerId, req.ContactId).First(&contactApply); res.Error != nil {
			if errors.Is(res.Error, gorm.ErrRecordNotFound) {
				contactApply = model.ContactApply{
					Uuid:        fmt.Sprintf("A%s", random.GetNowAndLenRandomString(11)),
					UserId:      req.OwnerId,
					ContactId:   req.ContactId,
					ContactType: contact_type_enum.User,
					Status:      contact_apply_status_enum.PENDING,
					Message:     req.Message,
					LastApplyAt: time.Now(),
				}
				if res := dao.GormDB.Create(&contactApply); res.Error != nil {
					zlog.Error(res.Error.Error())
					return "", res.Error
				}
			} else {
				zlog.Error(res.Error.Error())
				return "", res.Error
			}
		}
		contactApply.LastApplyAt = time.Now()

		if res := dao.GormDB.Save(&contactApply); res.Error != nil {
			zlog.Error(res.Error.Error())
			return "", res.Error
		}
		return "申请成功", nil
	} else if req.ContactId[0] == 'G' {
		var group model.GroupInfo
		if res := dao.GormDB.First(&group, "uuid = ?", req.ContactId); res.Error != nil {
			if errors.Is(res.Error, gorm.ErrRecordNotFound) {
				zlog.Info("contact not found in UserInfo")
				return "群聊不存在/已被禁用", nil
			} else {
				zlog.Error(res.Error.Error())
				return "", res.Error
			}
		}
		var contactApply model.ContactApply
		if res := dao.GormDB.Where("group_id = ? AND contact_id = ?", req.OwnerId, req.ContactId).First(&contactApply); res.Error != nil {
			if errors.Is(res.Error, gorm.ErrRecordNotFound) {
				contactApply = model.ContactApply{
					Uuid:        fmt.Sprintf("A%s", random.GetNowAndLenRandomString(11)),
					UserId:      req.OwnerId,
					ContactId:   req.ContactId,
					ContactType: contact_type_enum.Group,
					Status:      contact_apply_status_enum.PENDING,
					Message:     req.Message,
					LastApplyAt: time.Now(),
				}
				if res := dao.GormDB.Create(&contactApply); res.Error != nil {
					zlog.Error(res.Error.Error())
					return "", res.Error
				}
			} else {
				zlog.Error(res.Error.Error())
				return "", res.Error
			}
		}
		contactApply.LastApplyAt = time.Now()

		if res := dao.GormDB.Save(&contactApply); res.Error != nil {
			zlog.Error(res.Error.Error())
			return "", res.Error
		}
		return "申请成功", nil
	} else {
		return "用户/群聊不存在", nil
	}

}

// GetNewContactList 获取新的联系人申请列表
func (u *userContactService) GetNewContactList(ownerId string) (string, []respond.NewContactListRespond, error) {
	var contactApplyList []model.ContactApply
	if res := dao.GormDB.Where("contact_id = ? AND status = ?", ownerId, false).Find(&contactApplyList); res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			zlog.Info("没有在申请的联系人")
			return "没有在申请的联系人", nil, nil
		} else {
			zlog.Error(res.Error.Error())
			return "", nil, res.Error
		}
	}
	var rsp []respond.NewContactListRespond
	// 所有contact都没被禁用，被禁用的已经删除了
	for _, contactApply := range contactApplyList {
		var message string
		if contactApply.Message == "" {
			message = "申请理由：无"
		} else {
			message = "申请理由：" + contactApply.Message
		}
		newContact := respond.NewContactListRespond{
			ContactId: contactApply.Uuid,
			Message:   message,
		}
		var user model.UserInfo
		if res := dao.GormDB.First(&user, "uuid = ?", contactApply.UserId); res.Error != nil {
			return "", nil, res.Error
		}
		newContact.ContactId = user.Uuid
		newContact.ContactName = user.Nickname
		newContact.ContactAvatar = user.Avatar
		rsp = append(rsp, newContact)
	}
	return "获取成功", rsp, nil
}

// PassContactApply 通过联系人申请
func (u *userContactService) PassContactApply(ownerId string, contactId string) (string, error) {
	var contactApply model.ContactApply
	if res := dao.GormDB.Where("contact_id = ? AND user_id = ?", ownerId, contactId).First(&contactApply); res.Error != nil {
		zlog.Error(res.Error.Error())
		return "", res.Error
	}
	if contactApply.DeletedAt.Valid {
		zlog.Info("联系人已被禁用")
		return "联系人已被禁用", nil
	}
	contactApply.Status = contact_apply_status_enum.AGREE
	if res := dao.GormDB.Save(&contactApply); res.Error != nil {
		zlog.Error(res.Error.Error())
		return "", res.Error
	}
	newContact := model.UserContact{
		UserId:      ownerId,
		ContactId:   contactId,
		ContactType: contact_type_enum.User,     // 用户
		Status:      contact_status_enum.NORMAL, // 正常
		CreatedAt:   time.Now(),
		UpdateAt:    time.Now(),
	}
	if res := dao.GormDB.Create(&newContact); res.Error != nil {
		zlog.Error(res.Error.Error())
		return "", res.Error
	}
	newContact.UserId = contactId
	newContact.ContactId = ownerId
	if res := dao.GormDB.Create(&newContact); res.Error != nil {
		zlog.Error(res.Error.Error())
		return "", res.Error
	}
	return "已添加该联系人", nil
}
