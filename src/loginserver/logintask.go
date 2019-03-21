package main

import (
	"base/gnet"
	"base/log"
	"command"
	"github.com/golang/protobuf/proto"
)

type LoginTask struct {
	gnet.TCPTask

	msgHandler gnet.MessageHandler
}

func NewLoginTask() *LoginTask {
	task := &LoginTask{}
	task.init()

	return task
}

func (task *LoginTask) init() {
	task.Derived = task

	task.msgHandler.Reg(&command.ReqUserLogin{}, task.onReqUserLogin)
}

func (task *LoginTask) VerifyConn(msg *command.Message) bool {
	return true
}

func (task *LoginTask) RecycleConn() bool {
	return true
}

func (task *LoginTask) MsgParse(msg *command.Message) bool {

	task.msgHandler.Process(msg)
	return true
}

func (task *LoginTask) onReqUserLogin(cmd proto.Message) {

	msg := cmd.(*command.ReqUserLogin)

	log.Println("[登陆]收到登陆请求:", msg.Loginstr)

	gateway := gatewayManager.GetOne()
	if gateway == nil {
		log.Println("[登陆]没有可分配的网关节点")
	}

	snd := new(command.RetUserLogin)
	snd.Ip = gateway.Info.Ip
	snd.Port = gateway.Info.Port
	snd.Session = msg.Loginstr

	task.SendCmd(snd)
}
