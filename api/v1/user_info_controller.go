package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"kama_chat_server/internal/dto/request"
	"kama_chat_server/internal/service/gorm"
	"kama_chat_server/pkg/enum/error_info"
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
			"message": error_info.SYSTEM_ERROR,
		})
		return
	}
	fmt.Println(registerReq)
	message, userInfoStr, ret := gorm.UserInfoService.Register(c, registerReq)
	JsonBack(c, message, ret, userInfoStr)
}

// Login 登录
func Login(c *gin.Context) {
	var loginReq request.LoginRequest
	if err := c.BindJSON(&loginReq); err != nil {
		zlog.Error(err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": error_info.SYSTEM_ERROR,
		})
		return
	}
	message, userInfoStr, ret := gorm.UserInfoService.Login(c, loginReq)
	JsonBack(c, message, ret, userInfoStr)
}
