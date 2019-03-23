package gnet

import (
	"base/log"
	"base/util"
	"command"
	"github.com/golang/protobuf/proto"
	"reflect"
	"sync"
)

type MsgFunc func(proto.Message)

type MessageHandlerInfo struct {
	msgType    reflect.Type
	msgHandler MsgFunc
}

type MessageHandler struct {
	msgMap map[uint32]*MessageHandlerInfo

	mutex sync.Mutex
}

func (this *MessageHandler) Reg(msg proto.Message, fun MsgFunc) bool {

	name := proto.MessageName(msg)
	id := util.BKDRHash(name)

	info := new(MessageHandlerInfo)
	info.msgType = reflect.TypeOf(msg)
	info.msgHandler = fun

	this.mutex.Lock()
	defer this.mutex.Unlock()

	if this.msgMap == nil {
		this.msgMap = make(map[uint32]*MessageHandlerInfo)
	}

	this.msgMap[id] = info
	return true
}

func (this *MessageHandler) HaveMsgFunc(typeid uint32, name string) bool {

	_, ok := this.msgMap[typeid]
	if ok {
		return true
	}

	typeid = util.BKDRHash(name)
	_, ok = this.msgMap[typeid]

	return ok
}

func (this *MessageHandler) Process(msg *command.Message) bool {

	this.mutex.Lock()

	//log.Println("消息映射(", msg.Type, msg.Name, ")")
	info, ok := this.msgMap[msg.Type]
	if !ok {
		msg.Type = util.BKDRHash(msg.Name)
		info, ok = this.msgMap[msg.Type]
	}

	if !ok {

		log.Warnln("消息映射错误(", msg.Type, msg.Name, msg.Data, ")")
		return false
	}

	info = this.msgMap[msg.Type]
	cmd := reflect.New(info.msgType.Elem()).Interface()
	this.mutex.Unlock()

	if err := proto.Unmarshal(msg.Data, cmd.(proto.Message)); err != nil {
		log.Errorln("消息解析错误(", msg.Type, msg.Name, msg.Data, ")")
		return true
	}

	info.msgHandler(cmd.(proto.Message))

	return true
}
