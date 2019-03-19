package main

import (
	"base/gnet"
	"command"
)

type LoginTask struct {
	gnet.TCPTask
}

func NewLoginTask() *LoginTask {
	task := &LoginTask{}
	return task
}

func (task *LoginTask) VerifyConn(msg *command.Message) bool {
	return true
}

func (task *LoginTask) RecycleConn() bool {
	return true
}

func (task *LoginTask) MsgParse(msg *command.Message) bool {
	return true
}
