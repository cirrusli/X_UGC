package websocket

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"sync"
)

// Manager 所有websocket 信息
type Manager struct {
	//client.id => Client
	Group                map[string]*Client
	Lock                 sync.Mutex
	Register, UnRegister chan *Client
	Message              chan *MessageData
	clientCount          uint //分组及客户端数量
}

// WsManager 初始化 Manager 管理器
var WsManager = Manager{
	Group:       make(map[string]*Client),
	Register:    make(chan *Client, 128),
	UnRegister:  make(chan *Client, 128),
	Message:     make(chan *MessageData, 128),
	clientCount: 0,
}

// RegisterClient 注册客户端
func (manager *Manager) RegisterClient(client *Client) {
	manager.Register <- client
}

// UnRegisterClient 注销
func (manager *Manager) UnRegisterClient(client *Client) {
	manager.UnRegister <- client
}

// Start 启动websocket管理器
func (manager *Manager) Start() {
	slog.Info("websocket 服务器启动")
	for {
		select {
		case client := <-manager.Register:
			//注册客户端
			manager.Lock.Lock()
			manager.Group[client.Id] = client
			manager.clientCount += 1
			slog.Info(fmt.Sprintf("客户端注册: 客户端id为%s", client.Id))
			manager.Lock.Unlock()
		case client := <-manager.UnRegister:
			//注销客户端
			manager.Lock.Lock()
			if _, ok := manager.Group[client.Id]; ok {
				//关闭消息通道
				close(client.Message)
				//删除分组中客户
				delete(manager.Group, client.Id)
				//客户端数量减1
				manager.clientCount -= 1
				slog.Info(fmt.Sprintf("客户端注销: 客户端id为%s", client.Id))
			}
			manager.Lock.Unlock()
		case data := <-manager.Message:
			replyMsg, err := json.Marshal(&ReplyMsg{
				MsgType:  data.Message.Type,
				FromUser: data.FromUser,
				SendTime: data.Message.SendTime,
				Content:  data.Message.Content,
			})
			if err != nil {
				slog.Info("Marshal err=", err.Error())
			}
			if data.Message.Type == ToAllOnlineUser {
				//将数据广播给所有客户端
				for _, conn := range manager.Group {
					conn.Message <- replyMsg
				}
			} else { //指定用户
				for _, clientId := range data.Message.ClientIDs {
					if conn, ok := manager.Group[clientId]; ok {
						conn.Message <- replyMsg
					}
				}
			}
		}
	}
}
