package core

import (
	"github.com/gorilla/websocket"
	chatlog "logger/log"
	"time"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10
)

type Client struct {
	Id       int64
	Conn     *websocket.Conn
	SendPipe chan []byte
}

func NewClient(id int64, conn *websocket.Conn) *Client {
	return &Client{Id: id, Conn: conn,
		SendPipe: make(chan []byte, 512),
	}
}

//读取消息并执行对应的消息路由，有心跳检查的ping pong机制，每次pong之后都重新刷新长连接的存活时间
func (c *Client) Read() {
	defer func() {
		Manager.Unregister <- c
		_ = c.Conn.Close()
		close(c.SendPipe)
	}()

	_ = c.Conn.SetReadDeadline(time.Now().Add(pongWait))                                                           //设置读操作等待时间
	c.Conn.SetPongHandler(func(string) error { _ = c.Conn.SetReadDeadline(time.Now().Add(pongWait)); return nil }) // 每次pong结束后重新刷新读操作的等待时间

	for {
		mt, data, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				chatlog.Lg().Errorf("消息读取：网络错误: %v \n", err) //对端关闭时或特殊情况
			} else {
				chatlog.Lg().Errorf("pong超时 ：%v", err) //1分钟内对方未发送ping
			}
			break
		}
		if mt != websocket.BinaryMessage {
			chatlog.Lg().Errorln("消息类型读取出错")
			break
		}

		message := &MessageBase{}

		err = FromProtoc[*MessageBase](data, message)

		if err != nil {
			chatlog.Lg().Errorln("消息序列化出错")
			break
		}

		chatlog.Lg().Errorf("----------开始路由处理消息类型: %v", message.Impl.MsgType)

		err = Router.ExecHandler(&Context{
			Client: c,
			Msg:    message,
		})
		if err != nil {
			chatlog.Lg().Errorf("消息路由处理出错： %v", err)
			break
		}
	}
}

func (c *Client) Write() {
	//由于一个连接可以通过Read读取错误后进行关闭，所以Write没必要重复此操作
	ticker := time.NewTimer(pingPeriod)
	defer func() {
		ticker.Stop()
		_ = c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.SendPipe:
			_ = c.Conn.SetWriteDeadline(time.Now().Add(writeWait)) //设置写操作的等待时间
			if !ok {                                               //channel已经被关闭
				chatlog.Lg().Errorln("管道已经被关闭")
				return
			}

			w, err := c.Conn.NextWriter(websocket.BinaryMessage) //每次写消息都用新的Writer,会将之前的消息flush
			if err != nil {
				chatlog.Lg().Errorln("获取新Writer失败")
				return
			}

			w.Write(message)

			// 将还有在排队的消息继续写入
			n := len(c.SendPipe)
			for i := 0; i < n; i++ {
				w.Write(<-c.SendPipe)
			}

			if err := w.Close(); err != nil {
				chatlog.Lg().Errorf("Writer发生错误 %v", err)
				return
			}
		case <-ticker.C: //每过一个ping的等待时间重新刷新Writer的存活周期
			_ = c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil { //发送ping message
				chatlog.Lg().Errorln("ping 验证失败")
				return
			}
		}
	}
}
