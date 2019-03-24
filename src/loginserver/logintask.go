package main

import (
	"base/gnet"
	"base/log"
	"command"
	"github.com/golang/protobuf/proto"
)

var tempid uint32 = 1

type LoginTask struct {
	gnet.TCPTask

	tempid uint32

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

	cmd := new(command.ReqUserVerify)
	if err := proto.Unmarshal(msg.Data, cmd); err != nil {
		log.Println("[登陆]校验失败", task.GetIp())
		return false
	}

	loginTaskManager.UniqueAdd(task)

	snd := new(command.RetUserVerify)
	task.SendCmd(snd)

	log.Println("[登陆]登录校验", task.GetIp())
	return true
}

func (task *LoginTask) RecycleConn() bool {

	loginTaskManager.UniqueRemove(task)

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

	log.Println("[登陆]分配网关", gateway.Info.Ip, ":", gateway.Info.Port)
	snd := new(command.RetUserLogin)
	snd.Ip = gateway.Info.Ip
	snd.Port = gateway.Info.Port
	snd.Session = msg.Loginstr

	task.SendCmd(snd)
}

func (task *LoginTask) GetTempID() uint32 {

	if task.tempid == 0 {
		task.tempid = tempid
		tempid++
	}

	return task.tempid
}
