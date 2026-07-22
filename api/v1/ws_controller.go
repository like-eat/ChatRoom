package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"kama_chat_server/internal/dao"
	"kama_chat_server/internal/dto/request"
	"kama_chat_server/internal/model"
	"kama_chat_server/internal/service/auth"
	"kama_chat_server/internal/service/chat"
	"kama_chat_server/pkg/constants"
	"kama_chat_server/pkg/zlog"
	"net/http"
)

// WsLogin wss登录 Get
func WsLogin(c *gin.Context) {
	var tokenString string
	for _, protocol := range websocket.Subprotocols(c.Request) {
		if protocol != "kama-chat" {
			tokenString = protocol
			break
		}
	}
	claims, err := auth.ParseToken(tokenString)
	if err != nil {
		zlog.Error(err.Error())
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    http.StatusUnauthorized,
			"message": "WebSocket 登录令牌无效或已过期",
		})
		return
	}

	var user model.UserInfo
	if err := dao.GormDB.Select("uuid", "nickname", "avatar", "status").Where("uuid = ? AND status = 0", claims.Subject).First(&user).Error; err != nil {
		zlog.Error(err.Error())
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    http.StatusUnauthorized,
			"message": "用户不存在或账号已被禁用",
		})
		return
	}
	chat.NewClientInit(c, user.Uuid, user.Nickname, user.Avatar)
}

// WsLogout wss登出
func WsLogout(c *gin.Context) {
	var req request.WsLogoutRequest
	if err := c.BindJSON(&req); err != nil {
		zlog.Error(err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": constants.SYSTEM_ERROR,
		})
		return
	}
	message, ret := chat.ClientLogout(req.OwnerId)
	JsonBack(c, message, ret, nil)
}
