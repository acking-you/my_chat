package tests

import (
	"chat_socket/core"
	"chat_socket/serializer/models/protoc_message"
	"fmt"
	"testing"
)

func TestKConstants(t *testing.T) {
	fmt.Println(core.KText)
	ntc := &core.Notice{Impl: &protoc_message.Notice{
		StatusCode: 0,
		StatusMsg:  "safdsfsdfdsfafl",
	}}
	base := &core.MessageBase{Impl: &protoc_message.BaseMessage{
		MsgType:    -1,
		MsgContent: ntc.ToProtoc(),
	}}
	ntc2 := new(core.Notice)
	err := core.FromProtoc[*core.Notice](base.Impl.MsgContent, ntc2)
	if err != nil {
		panic(err)
	}
	fmt.Println(*ntc2)
}
