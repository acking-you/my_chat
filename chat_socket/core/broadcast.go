package core

import (
	"errors"
)

// IBroadcast 定义可广播的消息接口
type IBroadcast interface {
	MsgType() int32  //消息类型
	MsgRange() int32 //消息范围
	Sender() int64   //消息发送者
	Receiver() int64 //消息接收者
	Message() []byte //消息序列化后的产物
}

type Broadcast struct {
	Client  *Client    //源用户
	Msg     IBroadcast //所有广播的源信息
	clients *[]*Client //用于广播的目标
}

func (b *Broadcast) Do() error {
	var err error
	if err = b.getClients(); err != nil { //获取广播的目标用户
		return err
	}
	b.writeMessage() //向目标写入消息，具体来讲是写入到对方的管道中，故责任链已经转移

	return nil
}

// 获取目标用户和该消息的范围有关
func (b *Broadcast) getClients() error {
	//暂时只支持一对一
	var clients []*Client
	switch b.Msg.MsgRange() {
	case KOneToOne:
		//自己不能发信息给自己
		if b.Msg.Sender() == b.Msg.Receiver() {
			return errors.New("未定义行为")
		}
		userTo := b.Msg.Receiver()
		if cli, ok := Manager.ClientsMapper[userTo]; ok {
			clients = append(clients, cli)
			b.clients = &clients
			return nil
		}
		return errors.New("对方离线")
	default:
		return errors.New("不支持的消息类型")
	}
}

func (b *Broadcast) writeMessage() {
	//得到序列化消息，并放入对应的sendPipe
	for _, cli := range *b.clients {
		cli.SendPipe <- Reply(b.Msg.Message(), b.Msg.MsgType())
	}
}
