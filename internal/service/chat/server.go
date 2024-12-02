package chat

import (
	"apylee_chat_server/internal/model"
	"apylee_chat_server/pkg/zlog"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/gorilla/websocket"
)

type Server struct {
	Clients  map[string]*Client
	mutex    *sync.Mutex
	Transmit chan []byte  // 转发通道
	Login    chan *Client // 登录通道
	Logout   chan *Client // 推出登录通道
}

var ChatServer *Server

func InitServer() {
	if ChatServer == nil {
		ChatServer = &Server{
			Clients:  make(map[string]*Client),
			mutex:    &sync.Mutex{},
			Transmit: make(chan []byte),
			Login:    make(chan *Client),
			Logout:   make(chan *Client),
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
				s.Clients[client.Id] = client
				s.mutex.Unlock()
				err := client.Conn.WriteMessage(websocket.BinaryMessage, []byte("welcome to chat server"))
				if err != nil {
					zlog.Error(err.Error())
				}
				zlog.Info(fmt.Sprintf("Client %s logged in\n", client.Id))
			}

		case client := <-s.Logout:
			{
				s.mutex.Lock()
				delete(s.Clients, client.Id)
				s.mutex.Unlock()
				zlog.Info(fmt.Sprintf("Client %s logged out\n", client.Id))
			}

		case data := <-s.Transmit:
			{
				var message model.Message
				if err := json.Unmarshal(data, &message); err != nil {
					zlog.Error(err.Error())
				}
				if message.ReceiveId[0] == 'U' { // 发送给User
					receiveClient := s.Clients[message.ReceiveId]
					receiveClient.Send <- data // 向client.Send发送
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
							client.Send <- sendData // 可以使用Kafka
						}

					}
				}
			}
		}
	}
}
