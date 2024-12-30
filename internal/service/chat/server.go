package chat

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"kama_chat_server/internal/model"
	"kama_chat_server/pkg/constants"
	"kama_chat_server/pkg/zlog"
	"sync"
)

type Server struct {
	Clients  map[string]*Client
	mutex    *sync.Mutex
	Transmit chan []byte  // 转发通道
	Login    chan *Client // 登录通道
	Logout   chan *Client // 推出登录通道
}

var ChatServer *Server

func init() {
	if ChatServer == nil {
		ChatServer = &Server{
			Clients:  make(map[string]*Client),
			mutex:    &sync.Mutex{},
			Transmit: make(chan []byte, constants.CHANNEL_SIZE),
			Login:    make(chan *Client, constants.CHANNEL_SIZE),
			Logout:   make(chan *Client, constants.CHANNEL_SIZE),
		}
	}
}

// Start 启动函数，Server端用主进程起，Client端可以用协程起
func (s *Server) Start() {
	defer func() {
		close(s.Transmit)
		close(s.Logout)
		close(s.Login)
	}()
	for {
		select {
		case client := <-s.Login:
			{
				s.mutex.Lock()
				s.Clients[client.Uuid] = client
				s.mutex.Unlock()
				zlog.Debug(fmt.Sprintf("欢迎来到kama聊天服务器，亲爱的用户%s\n", client.Uuid))
				err := client.Conn.WriteMessage(websocket.TextMessage, []byte("欢迎来到kama聊天服务器"))
				if err != nil {
					zlog.Error(err.Error())
				}
			}

		case client := <-s.Logout:
			{
				s.mutex.Lock()
				delete(s.Clients, client.Uuid)
				s.mutex.Unlock()
				zlog.Info(fmt.Sprintf("Client %s logged out\n", client.Uuid))
			}

		case data := <-s.Transmit:
			{
				var message model.Message
				if err := json.Unmarshal(data, &message); err != nil {
					zlog.Error(err.Error())
				}
				if message.ReceiveId[0] == 'U' { // 发送给User
					receiveClient := s.Clients[message.ReceiveId]
					receiveClient.SendBack <- data // 向client.Send发送
				} else if message.ReceiveId[0] == 'G' { // 发送给Group
					var members []model.UserInfo
					for _, member := range members {
						// message是sendUserId给GroupId发送的
						// sendMessage是GroupId给除sendUserId之外的群成员发送的
						if member.Uuid != message.SendId {
							sendMessage := message
							sendMessage.SendId = message.ReceiveId
							sendMessage.ReceiveId = member.Uuid
							sendData, err := json.Marshal(sendMessage)
							if err != nil {
								zlog.Error(err.Error())
								break // 不转发了，直接结束
							}
							client := s.Clients[member.Uuid]
							client.SendBack <- sendData // 可以使用Kafka
						}

					}
				}
			}
		}
	}
}

func (s *Server) SendClientToLogin(client *Client) {
	s.mutex.Lock()
	s.Login <- client
	s.mutex.Unlock()
}

func (s *Server) SendClientToLogout(client *Client) {
	s.mutex.Lock()
	s.Logout <- client
	s.mutex.Unlock()
}

func (s *Server) SendMessageToTransmit(message []byte) {
	s.mutex.Lock()
	s.Transmit <- message
	s.mutex.Unlock()
}
