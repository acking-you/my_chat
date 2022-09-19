package main

import (
	wsrouter "chat_socket/router"
	"conf"
	"fmt"
	"go_http/router"
)

func main() {
	go wsrouter.Start() //启动websocket服务

	//启动http服务
	r := router.InitRouter()
	err := r.Run(fmt.Sprintf("%s:%d", conf.Config.Server.Host, conf.Config.Server.Port))
	if err != nil {
		return
	}
}
