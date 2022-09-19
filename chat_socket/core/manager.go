package core

import (
	chatlog "logger/log"
)

type ClientManager struct {
	ClientsMapper map[int64]*Client
	Broadcast     chan *Broadcast
	Register      chan *Client
	Unregister    chan *Client
}

var Manager = ClientManager{
	ClientsMapper: make(map[int64]*Client),
	Broadcast:     make(chan *Broadcast),
	Register:      make(chan *Client),
	Unregister:    make(chan *Client),
}

func (m *ClientManager) GetClient(userId int64) *Client {
	client := m.ClientsMapper[userId]
	return client
}

func (m *ClientManager) Do() {
	for {
		chatlog.Lg().Infoln("-------------------------websocket 开始监听信道---------------------------")
		select {
		case conn := <-Manager.Register: //建立连接注册用于管理
			m.register(conn)
		case conn := <-Manager.Unregister: //断开连接
			m.unregister(conn)
		case broadcast := <-Manager.Broadcast: //广播消息
			err := broadcast.Do()
			if err != nil {
				chatlog.Lg().Infof("消息广播未成功 clientId:%v errorInfo: %v", broadcast.Client.Id, err)
				broadcast.Client.SendPipe <- BadReplyMsg(err.Error())
			}
		}
	}
}

func (m *ClientManager) register(conn *Client) {
	chatlog.Lg().Infof("连接建立 clientId:%v addr:%v\n", conn.Id, conn.Conn.RemoteAddr())
	Manager.ClientsMapper[conn.Id] = conn
	conn.SendPipe <- GoodReplyMsg("连接服务器成功")
}

func (m *ClientManager) unregister(conn *Client) {
	chatlog.Lg().Infof("连接断开 clientId:%v addr:%v\n", conn.Id, conn.Conn.RemoteAddr())
	if _, ok := Manager.ClientsMapper[conn.Id]; ok {
		delete(Manager.ClientsMapper, conn.Id)
	}
}
