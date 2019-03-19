package gnet

import (
	"base/log"
	"base/util"
	"command"
	"github.com/golang/protobuf/proto"
)

type UserMsgFunc func(uint64, proto.Message) bool

type UserMsgParse struct {
	msgFunc map[uint32]UserMsgFunc
	typeMap map[uint32]proto.Message
}

func NewUserMsgParse() *UserMsgParse {
	msg := &UserMsgParse{
		msgFunc: make(map[uint32]UserMsgFunc),
		typeMap: make(map[uint32]proto.Message),
	}

	return msg
}

func (this *UserMsgParse) Reg(msg proto.Message, fun UserMsgFunc) bool {

	name := proto.MessageName(msg)
	id := util.BKDRHash(name)

	this.msgFunc[id] = fun
	this.typeMap[id] = msg
	return true
}

func (this *UserMsgParse) HaveMsgFunc(typeid uint32, name string) bool {

	_, ok := this.msgFunc[typeid]
	if ok {
		return true
	}

	typeid = util.BKDRHash(name)
	_, ok = this.msgFunc[typeid]

	return ok
}

func (this *UserMsgParse) Process(userid uint64, msg *command.Message) bool {

	fun, ok := this.msgFunc[msg.Type]
	if !ok {
		msg.Type = util.BKDRHash(msg.Name)
		fun, ok = this.msgFunc[msg.Type]
	}

	if !ok {
		return false
	}

	cmd := this.typeMap[msg.Type]
	if err := proto.Unmarshal(msg.Data, cmd); err != nil {
		log.Println("消息解析错误(", msg.Type, msg.Name, msg.Data, ")")
		return true
	}

	fun(userid, cmd)

	return true
}
