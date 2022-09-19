package router

import (
	"chat_socket/core"
	"chat_socket/handler"
)

func Start() {

	//websocket服务的路由
	core.Router.AddHandler(core.KText, handler.TextHandler)

	core.Manager.Do() //开启websocket服务的监听
}
