package chat

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"kama_chat_server/config"
	"kama_chat_server/pkg/constants"
	"kama_chat_server/pkg/zlog"
	"net/http"
)

type Client struct {
	Conn     *websocket.Conn
	Uuid     string
	SendTo   chan []byte // 给server端
	SendBack chan []byte // 给前端
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  2048,
	WriteBufferSize: 2048,
	// 检查连接的Origin头
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var messageMode = config.GetConfig().KafkaConfig.MessageMode

// 读取websocket消息并发送给send通道
func (c *Client) Read() {
	defer func() {
		if err := c.Conn.Close(); err != nil {
			zlog.Error(err.Error())
		}
		close(c.SendTo)
		close(c.SendBack)
	}()
	zlog.Info("ws read goroutine start")
	for {
		select {
		default:
			// 阻塞有一定隐患，因为下面要处理缓冲的逻辑，但是可以先不做优化，问题不大
			_, message, err := c.Conn.ReadMessage() // 阻塞状态
			if err != nil {
				zlog.Error(err.Error())
			} else {
				zlog.Info(string(message))
				if messageMode == "channel" {
					// 如果server的转发channel没满，先把sendto中的给transmit
					for len(ChatServer.Transmit) < constants.CHANNEL_SIZE && len(c.SendTo) > 0 {
						sendToMessage := <-c.SendTo
						ChatServer.SendMessageToTransmit(sendToMessage)
					}
					// 如果server没满，sendto空了，直接给server的transmit
					if len(ChatServer.Transmit) < constants.CHANNEL_SIZE {
						ChatServer.SendMessageToTransmit(message)
					} else {
						// 如果server满了，直接塞sendto
						c.SendTo <- message
					}
					// 否则考虑加宽channel size，或者使用kafka
					if err := c.Conn.WriteMessage(websocket.TextMessage, []byte("由于目前同一时间过多用户发送消息，消息发送失败，请稍后重试")); err != nil {
						zlog.Error(err.Error())
					}
				}
			}
		}
	}
}

// 从send通道读取消息发送给websocket
func (c *Client) Write() {
	defer func() {
		if err := c.Conn.Close(); err != nil {
			zlog.Error(err.Error())
		}
	}()
	zlog.Info("ws write goroutine start")
	for message := range c.SendBack { // 阻塞状态
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

	client := &Client{
		Conn:     conn,
		Uuid:     clientId,
		SendTo:   make(chan []byte, constants.CHANNEL_SIZE),
		SendBack: make(chan []byte, constants.CHANNEL_SIZE),
	}
	ChatServer.Login <- client
	go client.Read()
	go client.Write()
	zlog.Info("ws连接成功")
}
