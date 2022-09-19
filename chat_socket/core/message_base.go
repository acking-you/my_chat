package core

import (
	"chat_socket/serializer/models/protoc_message"
	"github.com/golang/protobuf/proto"
	chatlog "logger/log"
)

///所有的类型枚举值

const kNotice = -1

const (
	KText = iota
	KImage
)

const (
	KOneToOne  = iota //一对一聊天
	KOneToMany        //群组消息类型
)

const (
	KStatusOk = 0
)

///统一规范的消息MessageBase和实现protoc序列化接口

// MessageBase 用于规范所有的消息收发
type MessageBase struct {
	Impl *protoc_message.BaseMessage
}

func NewMessage(msgType int32, msgContent []byte) *MessageBase {
	return &MessageBase{
		Impl: &protoc_message.BaseMessage{
			MsgType:    msgType,
			MsgContent: msgContent,
		},
	}
}

func (m *MessageBase) ToProtoc() []byte {
	msg, err := proto.Marshal(m.Impl)
	if err != nil {
		chatlog.Lg().Errorln("message encoding failed")
		return nil
	}
	return msg
}

func (m *MessageBase) FromProtoc(msg []byte) error {
	m.Impl = &protoc_message.BaseMessage{}
	return proto.Unmarshal(msg, m.Impl)
}

///方便收发所有消息的接口和函数

// IProtoc 用于规范protoc序列化
type IProtoc interface {
	ToProtoc() []byte
	FromProtoc([]byte) error
}

// FromProtoc 接受用于反序列化的字节数据和符合IProtoc的类型，功能是将msg数据赋值初始化
func FromProtoc[T IProtoc](data []byte, msg T) error {
	err := msg.FromProtoc(data)
	if err != nil {
		return err
	}
	return nil
}

// Reply 接受protoc序列化后的二进制和对应的数据类型，返回经过BaseMessage包装后的protoc数据
func Reply(msg []byte, msgType int32) []byte {
	return NewMessage(msgType, msg).ToProtoc()
}

// GoodReplyMsg 服务端回应客户端的消息
func GoodReplyMsg(msg string) []byte {
	notice := &Notice{
		Impl: &protoc_message.Notice{StatusCode: 0,
			StatusMsg: msg},
	}
	baseMsg := &MessageBase{
		Impl: &protoc_message.BaseMessage{
			MsgType:    kNotice,
			MsgContent: notice.ToProtoc(),
		},
	}

	return baseMsg.ToProtoc()
}

func BadReplyMsg(msg string) []byte {
	notice := &Notice{
		Impl: &protoc_message.Notice{
			StatusCode: 1,
			StatusMsg:  msg,
		},
	}
	baseMsg := &MessageBase{
		Impl: &protoc_message.BaseMessage{
			MsgType:    kNotice,
			MsgContent: notice.ToProtoc(),
		},
	}

	return baseMsg.ToProtoc()
}
