package v1

import (
	"github.com/gin-gonic/gin"
	"kama_chat_server/internal/dto/request"
	"kama_chat_server/internal/service/gorm"
	"net/http"
)

// CreateSession 创建会话
//func CreateSession(c *gin.Context) {
//	var createSessionReq request.CreateSessionRequest
//	if err := c.BindJSON(&createSessionReq); err != nil {
//		c.JSON(http.StatusOK, gin.H{
//			"code":  400,
//			"error": err.Error(),
//		})
//		return
//	}
//	if err := gorm.SessionService.CreateSession(createSessionReq); err != nil {
//		c.JSON(http.StatusOK, gin.H{
//			"code":  400,
//			"error": err.Error(),
//		})
//		return
//	}
//	c.JSON(http.StatusOK, gin.H{
//		"code":    200,
//		"message": "create session success",
//	})
//}

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
	sessionList, err := gorm.SessionService.GetUserSessionList(getUserSessionListReq.OwnerId)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":  400,
			"error": err.Error(),
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": "load my session success",
			"data":    sessionList,
		})
	}
}

// GetGroupList 获取用户群聊列表
func GetGroupList(c *gin.Context) {
	var getGroupListReq request.OwnlistRequest
	if err := c.BindJSON(&getGroupListReq); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":  400,
			"error": err.Error(),
		})
		return
	}
	groupList, err := gorm.SessionService.GetGroupSessionList(getGroupListReq.OwnerId)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":  400,
			"error": err.Error(),
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": "load my group success",
			"data":    groupList,
		})
	}
}
