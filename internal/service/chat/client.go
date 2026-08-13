package chat

import (
	"context"
	"encoding/json"
	"kama_chat_server/internal/config"
	"kama_chat_server/internal/dao"
	"kama_chat_server/internal/dto/request"
	"kama_chat_server/internal/model"
	myKafka "kama_chat_server/internal/service/kafka"
	"kama_chat_server/pkg/constants"
	"kama_chat_server/pkg/enum/message/message_status_enum"
	"kama_chat_server/pkg/zlog"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/segmentio/kafka-go"
)

type MessageBack struct {
	Message []byte
	Uuid    string
}

type Client struct {
	Conn     *websocket.Conn
	Uuid     string
	Nickname string
	Avatar   string
	SendTo   chan []byte       // 给server端
	SendBack chan *MessageBack // 给前端
}

// http升级成websocket的升级器
var upgrader = websocket.Upgrader{
	ReadBufferSize:  2048,
	WriteBufferSize: 2048,
	Subprotocols:    []string{"kama-chat"},
	// 检查连接的Origin头
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var ctx = context.Background()

var messageMode = config.GetConfig().KafkaConfig.MessageMode

// 读取websocket消息并发送给send通道
// 从websocket读消息，发送到channel或者kafka
func (c *Client) Read() {
	zlog.Info("ws read goroutine start")
	for {
		// 阻塞有一定隐患，因为下面要处理缓冲的逻辑，但是可以先不做优化，问题不大
		_, jsonMessage, err := c.Conn.ReadMessage() // 阻塞状态
		if err != nil {
			zlog.Error(err.Error())
			return // 直接断开websocket
		} else {
			// 解析json消息
			var message = request.ChatMessageRequest{}
			if err := json.Unmarshal(jsonMessage, &message); err != nil {
				zlog.Error(err.Error())
				continue
			}
			// The authenticated WebSocket connection is the source of truth for
			// sender identity. Never trust sender fields supplied by the browser.
			message.SendId = c.Uuid
			message.SendName = c.Nickname
			message.SendAvatar = c.Avatar
			jsonMessage, err = json.Marshal(message)
			if err != nil {
				zlog.Error(err.Error())
				continue
			}
			log.Println("接受到消息为: ", jsonMessage)

			// 判断把消息发到哪里
			if messageMode == "channel" {
				// 如果server的转发channel没满，先把sendto中的给transmit
				// 这里使用自家邮箱和总邮箱来处理冲突流量
				// 如果总邮箱和自家邮箱都没满，就先往总邮箱投
				for len(ChatServer.Transmit) < constants.CHANNEL_SIZE && len(c.SendTo) > 0 {
					sendToMessage := <-c.SendTo
					ChatServer.SendMessageToTransmit(sendToMessage)
				}
				// Transmit是总邮箱，SendTo是自家邮箱
				if len(ChatServer.Transmit) < constants.CHANNEL_SIZE {
					ChatServer.SendMessageToTransmit(jsonMessage)
				} else if len(c.SendTo) < constants.CHANNEL_SIZE {
					// 如果server满了，直接塞sendto
					c.SendTo <- jsonMessage
				} else {
					// 否则考虑加宽channel size，或者使用kafka
					if err := c.Conn.WriteMessage(websocket.TextMessage, []byte("由于目前同一时间过多用户发送消息，消息发送失败，请稍后重试")); err != nil {
						zlog.Error(err.Error())
					}
				}
			} else {
				if err := myKafka.KafkaService.ChatWriter.WriteMessages(ctx, kafka.Message{
					Key:   []byte(strconv.Itoa(config.GetConfig().KafkaConfig.Partition)),
					Value: jsonMessage,
				}); err != nil {
					zlog.Error(err.Error())
				}
				zlog.Info("已发送消息：" + string(jsonMessage))
			}
		}
	}
}

// 从channel读数据往websocket里面写数据推送到前端
func (c *Client) Write() {
	zlog.Info("ws write goroutine start")
	for messageBack := range c.SendBack { // 阻塞状态
		// 把消息发给前端
		err := c.Conn.WriteMessage(websocket.TextMessage, messageBack.Message)
		if err != nil {
			zlog.Error(err.Error())
			return // 直接断开websocket
		}
		// log.Println("已发送消息：", messageBack.Message)
		// 说明顺利发送，修改消息的状态为已发送
		if res := dao.GormDB.Model(&model.Message{}).Where("uuid = ?", messageBack.Uuid).Update("status", message_status_enum.Sent); res.Error != nil {
			zlog.Error(res.Error.Error())
		}
	}
}

// NewClientInit 当接受到前端有登录消息时，会调用该函数
func NewClientInit(c *gin.Context, clientId, nickname, avatar string) {
	// 读取配置决定消息走哪条路
	kafkaConfig := config.GetConfig().KafkaConfig
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		zlog.Error(err.Error())
		return
	}

	// 造一个Client对象，含两条缓冲带的通道
	client := &Client{
		Conn:     conn,
		Uuid:     clientId,
		Nickname: nickname,
		Avatar:   avatar,
		SendTo:   make(chan []byte, constants.CHANNEL_SIZE),
		SendBack: make(chan *MessageBack, constants.CHANNEL_SIZE),
	}

	// 选择消息走哪条路
	if kafkaConfig.MessageMode == "channel" {
		ChatServer.SendClientToLogin(client)
	} else {
		KafkaChatServer.SendClientToLogin(client)
	}


	go client.Read() // 等前端发来的消息
	go client.Write() // 等服务器发来消息
	zlog.Info("ws连接成功")
}

// ClientLogout 当接受到前端有登出消息时，会调用该函数
func ClientLogout(clientId string) (string, int) {
	// 读取配置
	kafkaConfig := config.GetConfig().KafkaConfig
	// 获得是哪个client要退出（读共享 map 前必须加锁，防止并发读写 panic）
	ChatServer.mutex.Lock()
	client := ChatServer.Clients[clientId]
	ChatServer.mutex.Unlock()
	if client != nil {
		// channel和kafka都有自己的退出方案
		if kafkaConfig.MessageMode == "channel" {
			ChatServer.SendClientToLogout(client)
		} else {
			KafkaChatServer.SendClientToLogout(client)
		}
		if err := client.Conn.Close(); err != nil {
			zlog.Error(err.Error())
			return constants.SYSTEM_ERROR, -1
		}

		// 告诉所有等待这个通道的人，以后不会有人发消息了
		close(client.SendTo)
		close(client.SendBack)
	}
	return "退出成功", 0
}
