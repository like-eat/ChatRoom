package v1

import (
	"github.com/gin-gonic/gin"
	"kama_chat_server/internal/dto/request"
	"kama_chat_server/internal/service/gorm"
	"kama_chat_server/pkg/constants"
	"kama_chat_server/pkg/zlog"
	"net/http"
)

// CreateGroup 创建群聊
func CreateGroup(c *gin.Context) {
	var createGroupReq request.CreateGroupRequest
	if err := c.BindJSON(&createGroupReq); err != nil {
		zlog.Error(err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": constants.SYSTEM_ERROR,
		})
		return
	}
	message, ret := gorm.GroupInfoService.CreateGroup(createGroupReq)
	JsonBack(c, message, ret, nil)
}

// LoadMyGroup 获取我创建的群聊
func LoadMyGroup(c *gin.Context) {
	var loadMyGroupReq request.OwnlistRequest
	if err := c.BindJSON(&loadMyGroupReq); err != nil {
		zlog.Error(err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": constants.SYSTEM_ERROR,
		})
		return
	}
	message, groupList, ret := gorm.GroupInfoService.LoadMyGroup(loadMyGroupReq.OwnerId)
	JsonBack(c, message, ret, groupList)
}

// CheckGroupAddMode 检查群聊加群方式
func CheckGroupAddMode(c *gin.Context) {
	var req request.CheckGroupAddModeRequest
	if err := c.BindJSON(&req); err != nil {
		zlog.Error(err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": constants.SYSTEM_ERROR,
		})
		return
	}
	message, addMode, ret := gorm.GroupInfoService.CheckGroupAddMode(req.GroupId)
	JsonBack(c, message, ret, addMode)
}
