package utils

import (
	"reflect"

	"github.com/golang/protobuf/proto"
)

var Msg_CommonError uint16 = 0
var Msg_HHRequest uint16 = 1
var Msg_HHResponse uint16 = 2
var Msg_OrderResponse uint16 = 3
var Msg_AppOrderResponse uint16 = 4

type MessageHandler func(msgid uint16, msg interface{})

type MessageInfo struct {
	msgType   reflect.Type
	msgHanler MessageHandler
}

var (
	msg_map = make(map[uint16]MessageInfo)
)

func RegisterMessage(msgid uint16, msg interface{}, handler MessageHandler) {
	var info MessageInfo
	info.msgType = reflect.TypeOf(msg.(proto.Message))
	info.msgHanler = handler
	msg_map[msgid] = info
}
