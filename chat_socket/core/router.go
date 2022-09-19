package core

import (
	"errors"
)

type Context struct {
	Client *Client
	Msg    *MessageBase
}

type HandlerFunc func(ctx *Context)

type WsRouter struct {
	router map[int32]HandlerFunc
}

var Router = WsRouter{router: make(map[int32]HandlerFunc)}

// AddHandler 添加路由
func (w *WsRouter) AddHandler(msgType int32, handler HandlerFunc) {
	w.router[msgType] = handler
}

// ExecHandler 执行路由对应的handler
func (w *WsRouter) ExecHandler(ctx *Context) error {
	if ctx == nil {
		return errors.New("空指针错误")
	}
	handler, ok := w.router[ctx.Msg.Impl.MsgType]
	if ok {
		handler(ctx)
		return nil
	}
	return errors.New("没有对应的handler")
}
