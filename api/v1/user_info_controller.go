package v1

import (
	"apylee_chat_server/internal/dto/request"
	"apylee_chat_server/internal/service/gorm"
	"github.com/gin-gonic/gin"
	"net/http"
)

// Register 注册
func Register(c *gin.Context) {
	var registerReq request.RegisterRequest
	if err := c.BindJSON(&registerReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	message, err := gorm.UserInfoService.Register(registerReq)
	if message != "" && err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": message})
		return
	} else if message == "" && err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	} else if message == "" && err != nil {
		c.JSON(http.StatusOK, gin.H{"message": "register success"})
	}
}

// Login 登录
func Login(c *gin.Context) {
	var loginReq request.LoginRequest
	if err := c.BindJSON(&loginReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	message, err := gorm.UserInfoService.Login(loginReq)
	if message != "" && err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": message})
		return
	} else if message == "" && err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	} else if message == "" && err == nil {
		c.JSON(http.StatusOK, gin.H{"message": "login success"})
		return
	}
}
