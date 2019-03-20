package main

import (
	"base/gnet"
	"command"
)

type GatewayTask struct {
	gnet.TCPTask

	account string
}

func NewGatewayTask() *GatewayTask {
	task := &GatewayTask{}
	return task
}

func (task *GatewayTask) VerifyConn(msg *command.Message) bool {
	return true
}

func (task *GatewayTask) RecycleConn() bool {
	return true
}

func (task *GatewayTask) MsgParse(msg *command.Message) bool {
	return true
}

func (task *GatewayTask) GetAccount() string {
	return task.account
}
