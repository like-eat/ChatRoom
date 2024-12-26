package v1

import (
	"github.com/gin-gonic/gin"
	"kama_chat_server/internal/dto/request"
	"kama_chat_server/internal/service/gorm"
	"net/http"
)

// OpenSession 打开会话
func OpenSession(c *gin.Context) {
	var openSessionReq request.OpenSessionRequest
	if err := c.BindJSON(&openSessionReq); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":  400,
			"error": err.Error(),
		})
		return
	}
	message, sessionId, ret := gorm.SessionService.OpenSession(openSessionReq)
	JsonBack(c, message, ret, sessionId)
}

// GetUserSessionList 获取用户会话列表
func GetUserSessionList(c *gin.Context) {
	var getUserSessionListReq request.OwnlistRequest
	if err := c.BindJSON(&getUserSessionListReq); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":  400,
			"error": err.Error(),
		})
		return
	}
	message, sessionList, ret := gorm.SessionService.GetUserSessionList(getUserSessionListReq.OwnerId)
	JsonBack(c, message, ret, sessionList)
}

// GetGroupSessionList 获取用户群聊列表
func GetGroupSessionList(c *gin.Context) {
	var getGroupListReq request.OwnlistRequest
	if err := c.BindJSON(&getGroupListReq); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":  400,
			"error": err.Error(),
		})
		return
	}
	message, groupList, ret := gorm.SessionService.GetGroupSessionList(getGroupListReq.OwnerId)
	JsonBack(c, message, ret, groupList)
}

// DeleteSession 删除会话
func DeleteSession(c *gin.Context) {
	var deleteSessionReq request.DeleteSessionRequest
	if err := c.BindJSON(&deleteSessionReq); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":  400,
			"error": err.Error(),
		})
		return
	}
	message, ret := gorm.SessionService.DeleteSession(deleteSessionReq.SessionId)
	JsonBack(c, message, ret, nil)
}
