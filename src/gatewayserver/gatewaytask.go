package main

import (
	"base/gnet"
	"base/log"
	"base/util"
	"command"
	"github.com/golang/protobuf/proto"
)

type GatewayTask struct {
	gnet.TCPTask

	account string

	msgHandler gnet.MessageHandler
}

func NewGatewayTask() *GatewayTask {
	task := &GatewayTask{}
	task.init()

	return task
}

func (task *GatewayTask) init() {
	task.Derived = task

	task.msgHandler.Reg(&command.TestBroadcastAll{}, task.onTestBroadcastAll)
}

func (task *GatewayTask) VerifyConn(cmd *command.Message) bool {

	msg := new(command.ReqGatewayLogin)
	if err := proto.Unmarshal(cmd.Data, msg); err != nil {
		log.Println("[登陆]用户连接网关验证失败，消息解析错误")
		return false
	}

	params := util.Params2Map(msg.Session)

	task.account = params["account"]

	if len(task.GetAccount()) == 0 {
		log.Println("[登陆]帐号不能为空，session:", msg.Session)
		return true
	}

	if !gatewayTaskManager.UniqueAdd(task) {
		log.Println("[登陆]用户唯一性验证失败", task.account)
		return false
	}

	log.Println("[登陆]", task.account, "连接网关验证成功")

	snd := new(command.RetGatewayLogin)
	snd.Retcode = 0
	task.SendCmd(snd)

	task.queryRecord()

	return true
}

func (task *GatewayTask) RecycleConn() bool {

	log.Println("[登陆]", task.account, "与服务器断开连接")
	gatewayTaskManager.UniqueRemove(task)

	return true
}

func (task *GatewayTask) MsgParse(msg *command.Message) bool {

	task.msgHandler.Process(msg)

	return true
}

func (task *GatewayTask) GetAccount() string {
	return task.account
}

func (task *GatewayTask) queryRecord() {

	msg := new(command.ReqGateRegUser)
	msg.Account = task.GetAccount()

	routeManager.Broadcast(command.RecordServer, 0, msg)
}

func (task *GatewayTask) onTestBroadcastAll(cmd proto.Message) {
	gatewayTaskManager.BroadcastAll(cmd)
}
