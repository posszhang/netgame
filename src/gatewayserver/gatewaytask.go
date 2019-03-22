package main

import (
	"base/gnet"
	"base/log"
	"command"
	"github.com/golang/protobuf/proto"
)

type GatewayTask struct {
	gnet.TCPTask

	account string
}

func NewGatewayTask() *GatewayTask {
	task := &GatewayTask{}
	task.init()

	return task
}

func (task *GatewayTask) init() {
	task.Derived = task
}

func (task *GatewayTask) VerifyConn(cmd *command.Message) bool {

	log.Println("verify")

	msg := new(command.ReqGatewayLogin)
	if err := proto.Unmarshal(cmd.Data, msg); err != nil {
		log.Println("[登陆]用户连接网关验证失败，消息解析错误")
		return false
	}

	log.Println("verify1")

	task.account = msg.Session

	if !gatewayTaskManager.UniqueAdd(task) {
		log.Println("[登陆]用户唯一性验证失败", task.account)
		return false
	}

	log.Println("verify2")

	log.Println("[登陆]", task.account, "连接网关验证成功")
	return true
}

func (task *GatewayTask) RecycleConn() bool {

	gatewayTaskManager.UniqueRemove(task)

	return true
}

func (task *GatewayTask) MsgParse(msg *command.Message) bool {
	return true
}

func (task *GatewayTask) GetAccount() string {
	return task.account
}
