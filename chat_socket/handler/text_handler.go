package handler

import (
	"chat_socket/core"
	"github.com/gorilla/websocket"
	chatlog "logger/log"
)

// TextHandler 文字类型的消息处理（只需广播即可
func TextHandler(c *core.Context) {
	data := c.Msg.Impl.MsgContent
	textMsg := new(core.TextMessage)
	err := core.FromProtoc(data, textMsg)
	chatlog.Lg().Errorf("序列化后的textMsg :%v \n", *textMsg)
	//发生解析错误
	if err != nil {
		_ = c.Client.Conn.WriteMessage(websocket.BinaryMessage, core.BadReplyMsg("消息类型与解析类型不对应"))
		return
	}

	//开始进行消息的广播
	core.Manager.Broadcast <- &core.Broadcast{
		Client: c.Client,
		Msg:    textMsg,
	}
}
