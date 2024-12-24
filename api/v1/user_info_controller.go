package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"kama_chat_server/internal/dto/request"
	"kama_chat_server/internal/service/gorm"
	"kama_chat_server/pkg/zlog"
	"net/http"
)

// Register 注册
func Register(c *gin.Context) {
	var registerReq request.RegisterRequest
	if err := c.BindJSON(&registerReq); err != nil {
		zlog.Error(err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "系统问题，请联系工作人员",
		})
		return
	}
	fmt.Println(registerReq)
	message, userInfoStr := gorm.UserInfoService.Register(c, registerReq)
	if message != "" {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": message,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": message,
			"data":    userInfoStr,
		})
	}
}

// Login 登录
func Login(c *gin.Context) {
	var loginReq request.LoginRequest
	if err := c.BindJSON(&loginReq); err != nil {
		zlog.Error(err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "系统问题，请联系工作人员",
		})
		return
	}
	message, userInfoStr := gorm.UserInfoService.Login(c, loginReq)
	if message != "" {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": message,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": message,
			"data":    userInfoStr,
		})
	}
}
