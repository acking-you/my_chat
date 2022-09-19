package core

import (
	"chat_socket/serializer/models/protoc_message"
	"github.com/golang/protobuf/proto"
	chatlog "logger/log"
)

type TextMessage struct {
	Impl *protoc_message.TextMessage
}

type Notice struct {
	Impl *protoc_message.Notice
}

/// 实现IProtoc接口

func (t *TextMessage) ToProtoc() []byte {
	ret, err := proto.Marshal(t.Impl)
	if err != nil {
		chatlog.Lg().Errorln("TextMessage encoding failed")
		return nil
	}
	//打印反序列化结果进行测试
	text := new(TextMessage)
	err = FromProtoc(ret, text)
	if err != nil {
		chatlog.Lg().Debugln("反序列化测试失败")
		return ret
	}
	chatlog.Lg().Debugf("反序列化测试 %v\n", *text)
	return ret
}

func (t *TextMessage) FromProtoc(msg []byte) error {
	t.Impl = &protoc_message.TextMessage{}
	return proto.Unmarshal(msg, t.Impl)
}

func (t *Notice) ToProtoc() []byte {
	ret, err := proto.Marshal(t.Impl)
	if err != nil {
		chatlog.Lg().Errorln("NoticeMessage encoding failed")
		return nil
	}
	return ret
}

func (t *Notice) FromProtoc(msg []byte) error {
	t.Impl = &protoc_message.Notice{}
	return proto.Unmarshal(msg, t.Impl)
}

///实现IReply接口

func (t *TextMessage) MsgType() int32 {
	return KText
}

///实现IBroadcast接口

func (t *TextMessage) Sender() int64 {
	return t.Impl.Sender
}

func (t *TextMessage) Receiver() int64 {
	return t.Impl.Receiver
}

func (t *TextMessage) MsgRange() int32 {
	return t.Impl.MsgRange
}

func (t *TextMessage) Message() []byte {
	return t.ToProtoc()
}
