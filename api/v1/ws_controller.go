package v1

import (
	"kama_chat_server/internal/dao"
	"kama_chat_server/internal/dto/request"
	"kama_chat_server/internal/model"
	"kama_chat_server/internal/service/auth"
	"kama_chat_server/internal/service/chat"
	"kama_chat_server/pkg/constants"
	"kama_chat_server/pkg/zlog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// go中没有单独编写WebSocket，因此要用第三方库
// wesocket.Upgrader把http升级成WebSocket

// 走REST API是一问一答
// 走WebSocket是持续连接，因此要用长连接

// WsLogin wss登录 Get
func WsLogin(c *gin.Context) {
	var tokenString string
	// 把token从请求里拿出来
	for _, protocol := range websocket.Subprotocols(c.Request) {
		if protocol != "kama-chat" {
			tokenString = protocol
			break
		}
	}

	// 验证token是否有效
	claims, err := auth.ParseToken(tokenString)
	if err != nil {
		zlog.Error(err.Error())
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    http.StatusUnauthorized,
			"message": "WebSocket 登录令牌无效或已过期",
		})
		return
	}

	// 查询用户是否还在数据库中
	var user model.UserInfo
	if err := dao.GormDB.Select("uuid", "nickname", "avatar", "status").Where("uuid = ? AND status = 0", claims.Subject).First(&user).Error; err != nil {
		zlog.Error(err.Error())
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    http.StatusUnauthorized,
			"message": "用户不存在或账号已被禁用",
		})
		return
	}

	// 创建一个客户端实例
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
	// 这个才是关键，真正的退出逻辑
	message, ret := chat.ClientLogout(req.OwnerId)
	JsonBack(c, message, ret, nil)
}
