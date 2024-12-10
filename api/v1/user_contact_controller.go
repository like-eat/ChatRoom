package v1

import (
	"github.com/gin-gonic/gin"
	"kama_chat_server/internal/dto/request"
	"kama_chat_server/internal/service/gorm"
	"net/http"
)

// GetUserList 获取联系人列表
func GetUserList(c *gin.Context) {
	var myUserListReq request.MyUserListRequest
	if err := c.BindJSON(&myUserListReq); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":  400,
			"error": err.Error(),
		})
	}
	message, userList, err := gorm.UserContactService.GetUserList(myUserListReq.OwnerId)
	if message == "" && err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":  400,
			"error": err.Error(),
		})
	} else if message != "" && err == nil {
		c.JSON(http.StatusOK, gin.H{
			"code":  400,
			"error": message,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": "get userlist success",
			"data":    userList,
		})
	}
}
