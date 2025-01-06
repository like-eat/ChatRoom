package v1

import (
	"github.com/gin-gonic/gin"
	"kama_chat_server/internal/service/chat"
	"kama_chat_server/pkg/zlog"
	"net/http"
)

// WsHandler 处理websocket请求
func WsHandler(c *gin.Context) {
	clientId := c.Query("client_id")
	if clientId == "" {
		zlog.Error("clientId获取失败")
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "clientId获取失败",
		})
		return
	}
	chat.NewClientInit(c, clientId)
}
