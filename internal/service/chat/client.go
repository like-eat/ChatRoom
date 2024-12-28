package chat

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"kama_chat_server/pkg/constants"
	"kama_chat_server/pkg/zlog"
	"net/http"
	"time"
)

type Client struct {
	Conn *websocket.Conn
	Uuid string
	Send chan []byte
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  2048,
	WriteBufferSize: 2048,
	// 检查连接的Origin头
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// 读取websocket消息并发送给send通道
func (c *Client) Read() {
	defer func() {
		if err := c.Conn.Close(); err != nil {
			zlog.Error(err.Error())
		}
		close(c.Send)
	}()
	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			zlog.Error(err.Error())
		} else {
			c.Send <- message
		}
		time.Sleep(1)
	}
}

// 从send通道读取消息发送给websocket
func (c *Client) Write() {
	defer func() {
		if err := c.Conn.Close(); err != nil {
			zlog.Error(err.Error())
		}
	}()
	for message := range c.Send {
		// 通过 WebSocket 发送消息
		err := c.Conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			zlog.Error(err.Error())
		}
	}
}

// NewClientInit 当接受到前端有登录消息时，会调用该函数
func NewClientInit(c *gin.Context, clientId string) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		zlog.Error(err.Error())
	}
	defer func() {
		if err := conn.Close(); err != nil {
			zlog.Error(err.Error())
		}
	}()

	client := &Client{
		Conn: conn,
		Uuid: clientId,
		Send: make(chan []byte, constants.CHANNEL_SIZE),
	}
	ChatServer.Login <- client
	//go client.Read()
	//go client.Write()
	zlog.Info("ws连接成功")
}
