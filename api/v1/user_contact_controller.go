package v1

import (
	"github.com/gin-gonic/gin"
	"kama_chat_server/internal/dto/request"
	"kama_chat_server/internal/service/gorm"
	"log"
	"net/http"
)

// GetUserList 获取联系人列表
func GetUserList(c *gin.Context) {
	var myUserListReq request.OwnlistRequest
	if err := c.BindJSON(&myUserListReq); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":  400,
			"error": err.Error(),
		})
	}
	message, userList, ret := gorm.UserContactService.GetUserList(myUserListReq.OwnerId)
	JsonBack(c, message, ret, userList)
}

// LoadMyJoinedGroup 获取我加入的群聊
func LoadMyJoinedGroup(c *gin.Context) {
	var loadMyJoinedGroupReq request.OwnlistRequest
	if err := c.BindJSON(&loadMyJoinedGroupReq); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":  400,
			"error": err.Error(),
		})
		return
	}
	message, groupList, ret := gorm.UserContactService.LoadMyJoinedGroup(loadMyJoinedGroupReq.OwnerId)
	JsonBack(c, message, ret, groupList)
}

// GetContactInfo 获取联系人信息
func GetContactInfo(c *gin.Context) {
	var getContactInfoReq request.GetContactInfoRequest
	if err := c.BindJSON(&getContactInfoReq); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":  400,
			"error": err.Error(),
		})
		return
	}
	log.Println(getContactInfoReq)
	message, contactInfo, ret := gorm.UserContactService.GetContactInfo(getContactInfoReq.ContactId)
	JsonBack(c, message, ret, contactInfo)
}

// DeleteContact 删除联系人
func DeleteContact(c *gin.Context) {
	var deleteContactReq request.DeleteContactRequest
	if err := c.BindJSON(&deleteContactReq); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":  400,
			"error": err.Error(),
		})
		return
	}
	message, ret := gorm.UserContactService.DeleteContact(deleteContactReq.OwnerId, deleteContactReq.ContactId)
	JsonBack(c, message, ret, nil)
}

// ApplyContact 申请添加联系人
func ApplyContact(c *gin.Context) {
	var applyContactReq request.ApplyContactRequest
	if err := c.BindJSON(&applyContactReq); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":  400,
			"error": err.Error(),
		})
		return
	}
	message, ret := gorm.UserContactService.ApplyContact(applyContactReq)
	JsonBack(c, message, ret, nil)
}

// GetNewContactList 获取新的联系人申请列表
func GetNewContactList(c *gin.Context) {
	var req request.OwnlistRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":  400,
			"error": err.Error(),
		})
		return
	}
	message, data, ret := gorm.UserContactService.GetNewContactList(req.OwnerId)
	JsonBack(c, message, ret, data)
}

// PassContactApply 通过联系人申请
func PassContactApply(c *gin.Context) {
	var passContactApplyReq request.PassContactApplyRequest
	if err := c.BindJSON(&passContactApplyReq); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":  400,
			"error": err.Error(),
		})
		return
	}
	message, ret := gorm.UserContactService.PassContactApply(passContactApplyReq.OwnerId, passContactApplyReq.ContactId)
	JsonBack(c, message, ret, nil)
}
