package api

import (
	"chat_socket/core"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go_http/pkg/utils"
	chatlog "logger/log"
	"net/http"
)

var upgrader = websocket.Upgrader{ReadBufferSize: 1024, WriteBufferSize: 1024, CheckOrigin: func(r *http.Request) bool {
	return true
}}

func WebSocketRegister(c *gin.Context) {
	//获取userId
	userId := utils.GetId(c)
	if userId == -1 {
		return
	}

	//升级为websocket协议
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		chatlog.Lg().Errorln(err)
		http.NotFound(c.Writer, c.Request)
		return
	}

	//包装websocket并注册管理
	client := core.NewClient(userId, conn)
	core.Manager.Register <- client

	go client.Read() //读写进行异步处理
	go client.Write()
}
